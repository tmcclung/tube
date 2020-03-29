package app

import (
	"encoding/json"
	"os"
)

// Config settings for main App.
type Config struct {
	Library     []*PathConfig      `json:"library"`
	Server      *ServerConfig      `json:"server"`
	Thumbnailer *ThumbnailerConfig `json:"thumbnailer"`
	Transcoder  *TranscoderConfig  `json:"transcoder"`
	Feed        *FeedConfig        `json:"feed"`
}

// PathConfig settings for media library path.
type PathConfig struct {
	Path   string `json:"path"`
	Prefix string `json:"prefix"`
}

// ServerConfig settings for App Server.
type ServerConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	StorePath     string `json:"store_path"`
	UploadPath    string `json:"upload_path"`
	MaxUploadSize int64  `json:"max_upload_size"`
}

// ThumbnailerConfig settings for Transcoder
type ThumbnailerConfig struct {
	Timeout int `json:"timeout"`
}

// TranscoderConfig settings for Transcoder
type TranscoderConfig struct {
	Timeout int `json:"timeout"`
}

// FeedConfig settings for App Feed.
type FeedConfig struct {
	ExternalURL string `json:"external_url"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Author      struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Copyright string `json:"copyright"`
}

// DefaultConfig returns Config initialized with default values.
func DefaultConfig() *Config {
	return &Config{
		Library: []*PathConfig{
			&PathConfig{
				Path:   "videos",
				Prefix: "",
			},
		},
		Server: &ServerConfig{
			Host:          "0.0.0.0",
			Port:          8000,
			StorePath:     "tube.db",
			UploadPath:    "uploads",
			MaxUploadSize: 104857600,
		},
		Thumbnailer: &ThumbnailerConfig{
			Timeout: 60,
		},
		Transcoder: &TranscoderConfig{
			Timeout: 300,
		},
		Feed: &FeedConfig{
			ExternalURL: "http://localhost:8000",
		},
	}
}

// ReadFile reads a JSON file into Config.
func (c *Config) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	return d.Decode(c)
}
