package importers

import (
	"errors"
	"strings"
)

var (
	ErrUnsupportedVideoURL = errors.New("error: unsupported video url")
)

type VideoInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`

	VideoURL     string `json:"video_url"`
	ThumbnailURL string `json:"thumbnail_url"`
}

type Importer interface {
	GetVideoInfo(url string) (VideoInfo, error)
}

func NewImporter(url string) (Importer, error) {
	if strings.Contains(url, "youtube.com") || strings.HasPrefix(strings.ToLower(url), "youtube:") {
		return &YoutubeImporter{}, nil
	} else if strings.Contains(url, "vimeo.com") || strings.HasPrefix(strings.ToLower(url), "vimeo:") {
		return &VimeoImporter{}, nil
	} else {
		return nil, ErrUnsupportedVideoURL
	}
}
