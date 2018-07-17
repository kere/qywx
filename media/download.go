package media

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	urlDownload = "https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)

// DownloadMedia 下载临时素材
func DownloadMedia(mediaID, token, name, folder string) (string, error) {
	uri := fmt.Sprintf(urlDownload, token, mediaID)

	resq, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	defer resq.Body.Close()
	src, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		return "", err
	}

	// fmt.Println(resq.Header.Get("Content-Type"))
	// fmt.Println(resq.Header.Get("Content-disposition"))
	// attachment; filename="856e35aa2fd8c8d7682614bff85fba19.png"
	tmp := resq.Header.Get("Content-disposition")
	oname := tmp[22 : len(tmp)-1]
	ext := filepath.Ext(oname)

	name += ext
	err = ioutil.WriteFile(filepath.Join(folder, name), src, os.ModePerm)

	return name, err
}
