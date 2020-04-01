package importers

import (
	"fmt"
	"strings"

	"github.com/rylio/ytdl"
)

type YoutubeImporter struct{}

func (i *YoutubeImporter) GetVideoInfo(url string) (videoInfo VideoInfo, err error) {
	if strings.HasPrefix(strings.ToLower(url), "youtube:") {
		url = strings.TrimSpace(strings.SplitN(url, ":", 2)[1])
	}

	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		err = fmt.Errorf("error retriving youtube video info: %w", err)
		return
	}

	videoURL, err := ytdl.DefaultClient.GetDownloadURL(info, info.Formats[0])
	if err != nil {
		err = fmt.Errorf("error retriving youtube video  url: %w", err)
		return
	}
	videoInfo.VideoURL = videoURL.String()

	videoInfo.ThumbnailURL = info.GetThumbnailURL(ytdl.ThumbnailQualityHigh).String()

	videoInfo.ID = info.ID
	videoInfo.Title = info.Title
	videoInfo.Description = info.Description

	return
}
