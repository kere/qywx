package media

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const (
	urlDownload = "https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)

// DownloadMedia 下载临时素材
func DownloadMedia(mediaID, token, name, folder string) (string, error) {
	uri := fmt.Sprintf(urlDownload, token, mediaID)
	imgSize := uint(800)
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
	tempDir := os.TempDir()
	tempName := filepath.Join(tempDir, name)

	err = ioutil.WriteFile(tempName, src, os.ModePerm)
	if err != nil {
		return "", err
	}

	file, err := os.Open(tempName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if strings.ToLower(ext) == ".jpg" {
		// decode jpeg into image.Image
		img, err1 := jpeg.Decode(file)
		if err1 != nil {
			return "", err1
		}

		m := resize.Resize(imgSize, 0, img, resize.Lanczos3)

		out, err1 := os.Create(filepath.Join(folder, name))
		if err1 != nil {
			return "", err1
		}
		defer out.Close()

		jpeg.Encode(out, m, nil)

	} else if strings.ToLower(ext) == ".png" {
		// decode jpeg into image.Image
		img, err1 := png.Decode(file)
		if err1 != nil {
			return "", err1
		}

		m := resize.Resize(imgSize, 0, img, resize.Lanczos3)

		out, err1 := os.Create(filepath.Join(folder, name))
		if err1 != nil {
			return "", err1
		}
		defer out.Close()

		png.Encode(out, m)
	}
	os.Remove(tempName)

	return name, err
}
