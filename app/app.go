// Package app manages main application server.
package app

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"github.com/renstrom/shortuuid"
	log "github.com/sirupsen/logrus"
	"github.com/wybiral/tube/media"
	"github.com/wybiral/tube/utils"
)

//go:generate rice embed-go

// App represents main application.
type App struct {
	Config    *Config
	Library   *media.Library
	Watcher   *fsnotify.Watcher
	Templates *templateStore
	Feed      []byte
	Listener  net.Listener
	Router    *mux.Router
}

// NewApp returns a new instance of App from Config.
func NewApp(cfg *Config) (*App, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	a := &App{
		Config: cfg,
	}
	// Setup Library
	a.Library = media.NewLibrary()
	// Setup Watcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	a.Watcher = w
	// Setup Listener
	ln, err := newListener(cfg.Server)
	if err != nil {
		return nil, err
	}
	a.Listener = ln

	// Templates
	box := rice.MustFindBox("../templates")

	a.Templates = newTemplateStore("base")

	indexTemplate := template.New("index")
	template.Must(indexTemplate.Parse(box.MustString("index.html")))
	template.Must(indexTemplate.Parse(box.MustString("base.html")))
	a.Templates.Add("index", indexTemplate)

	uploadTemplate := template.New("upload")
	template.Must(uploadTemplate.Parse(box.MustString("upload.html")))
	template.Must(uploadTemplate.Parse(box.MustString("base.html")))
	a.Templates.Add("upload", uploadTemplate)

	// Setup Router
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", a.indexHandler).Methods("GET")
	r.HandleFunc("/upload", a.uploadHandler).Methods("GET", "POST")
	r.HandleFunc("/v/{id}.mp4", a.videoHandler).Methods("GET")
	r.HandleFunc("/v/{prefix}/{id}.mp4", a.videoHandler).Methods("GET")
	r.HandleFunc("/t/{id}", a.thumbHandler).Methods("GET")
	r.HandleFunc("/t/{prefix}/{id}", a.thumbHandler).Methods("GET")
	r.HandleFunc("/v/{id}", a.pageHandler).Methods("GET")
	r.HandleFunc("/v/{prefix}/{id}", a.pageHandler).Methods("GET")
	r.HandleFunc("/feed.xml", a.rssHandler).Methods("GET")
	// Static file handler
	fsHandler := http.StripPrefix(
		"/static",
		http.FileServer(rice.MustFindBox("../static").HTTPBox()),
	)
	r.PathPrefix("/static/").Handler(fsHandler).Methods("GET")
	a.Router = r
	return a, nil
}

// Run imports the library and starts server.
func (a *App) Run() error {
	for _, pc := range a.Config.Library {
		p := &media.Path{
			Path:   pc.Path,
			Prefix: pc.Prefix,
		}
		err := a.Library.AddPath(p)
		if err != nil {
			return err
		}
		err = a.Library.Import(p)
		if err != nil {
			return err
		}
		a.Watcher.Add(p.Path)
	}
	buildFeed(a)
	go startWatcher(a)
	return http.Serve(a.Listener, a.Router)
}

func (a *App) render(name string, w http.ResponseWriter, ctx interface{}) {
	buf, err := a.Templates.Exec(name, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HTTP handler for /
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/")
	pl := a.Library.Playlist()
	if len(pl) > 0 {
		http.Redirect(w, r, "/v/"+pl[0].ID, 302)
	} else {
		ctx := &struct {
			Playing  *media.Video
			Playlist media.Playlist
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
		}

		a.render("index", w, ctx)
	}
}

// HTTP handler for /upload
func (a *App) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ctx := &struct{}{}
		a.render("upload", w, ctx)
	} else if r.Method == "POST" {
		// TODO: Move to a constant
		r.ParseMultipartForm((10 << 20) * 10) // 100MB

		file, handler, err := r.FormFile("video_file")
		if err != nil {
			err := fmt.Errorf("error processing form: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// TODO: Make collection user selectable from drop-down in Form
		// XXX: Assume we can put uploaded videos into the first collection (sorted) we find
		keys := make([]string, 0, len(a.Library.Paths))
		for k := range a.Library.Paths {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		collection := keys[0]

		fn := filepath.Join(
			a.Config.Server.UploadPath,
			fmt.Sprintf(
				"%s%s",
				shortuuid.New(),
				filepath.Ext(handler.Filename),
			),
		)

		f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			err := fmt.Errorf("error opening file for writing: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			err := fmt.Errorf("error writing file: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tf, err := ioutil.TempFile("", "tube-upload-*.mp4")
		if err != nil {
			err := fmt.Errorf("error creating tempory file: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer os.Remove(tf.Name())

		of := filepath.Join(
			a.Library.Paths[collection].Path,
			fmt.Sprintf(
				"%s.mp4",
				strings.TrimSuffix(filepath.Base(fn), filepath.Ext(fn)),
			),
		)

		// TODO: Use a proper Job Queue and make this async
		if err := utils.RunCmd(
			// TODO: Make thie timeout configurable?
			60,
			"ffmpeg",
			"-y",
			"-i", fn,
			"-vcodec", "h264",
			"-acodec", "aac",
			"-strict", "-2",
			"-loglevel", "quiet",
			tf.Name(),
		); err != nil {
			err := fmt.Errorf("error transcoding video: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		os.Rename(tf.Name(), of)
		os.Remove(fn)

		if err := a.Library.Add(of); err != nil {
			err := fmt.Errorf("error adding new video: %w", err)
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Video successfully uploaded!")
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// HTTP handler for /v/id
func (a *App) pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/v/%s", id)
	playing, ok := a.Library.Videos[id]
	if !ok {
		ctx := &struct {
			Playing  *media.Video
			Playlist media.Playlist
		}{
			Playing:  &media.Video{ID: ""},
			Playlist: a.Library.Playlist(),
		}
		a.render("upload", w, ctx)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := &struct {
		Playing  *media.Video
		Playlist media.Playlist
	}{
		Playing:  playing,
		Playlist: a.Library.Playlist(),
	}
	a.render("index", w, ctx)
}

// HTTP handler for /v/id.mp4
func (a *App) videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/v/%s", id)
	m, ok := a.Library.Videos[id]
	if !ok {
		return
	}
	title := m.Title
	disposition := "attachment; filename=\"" + title + ".mp4\""
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeFile(w, r, m.Path)
}

// HTTP handler for /t/id
func (a *App) thumbHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	prefix, ok := vars["prefix"]
	if ok {
		id = path.Join(prefix, id)
	}
	log.Printf("/t/%s", id)
	m, ok := a.Library.Videos[id]
	if !ok {
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	if m.ThumbType == "" {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(rice.MustFindBox("../static").MustBytes("defaulticon.jpg"))
	} else {
		w.Header().Set("Content-Type", m.ThumbType)
		w.Write(m.Thumb)
	}
}

// HTTP handler for /feed.xml
func (a *App) rssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	w.Header().Set("Content-Type", "text/xml")
	w.Write(a.Feed)
}
