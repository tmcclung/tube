package importers

import "fmt"

type VimeoImporter struct{}

func (i *VimeoImporter) GetVideoInfo(url string) (videoInfo VideoInfo, err error) {
	err = fmt.Errorf("Not Implemented")
	return
}
