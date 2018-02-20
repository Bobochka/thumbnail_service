package downloader

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/Bobochka/thumbnail_service/lib"
)

type Http struct {
}

var supportedTypes = map[string]struct{}{
	"image/jpeg": struct{}{},
	"image/png":  struct{}{},
	"image/gif":  struct{}{},
}

func (Http) Download(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, lib.NewError(err, lib.ResourceUnreachable)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, lib.NewError(err, lib.ResourceUnreachable)
	}

	t := http.DetectContentType(data)
	if _, ok := supportedTypes[t]; !ok {
		return nil, lib.NewError(fmt.Errorf("content type %s not supported"), lib.UnsupportedContentType)
	}

	return data, nil
}
