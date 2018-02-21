package downloader

import (
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/Bobochka/thumbnail_service/lib"
)

type Http struct {
	contentTypes map[string]struct{}
}

func New(allowedContentTypes []string) *Http {
	ct := map[string]struct{}{}
	for _, t := range allowedContentTypes {
		ct[t] = struct{}{}
	}

	return &Http{
		contentTypes: ct,
	}
}

func (d *Http) Download(url string) ([]byte, error) {
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

	if resp.StatusCode/100 != 2 {
		return nil, lib.NewError(err, lib.ResourceUnreachable)
	}

	t := http.DetectContentType(data)
	if _, ok := d.contentTypes[t]; !ok {
		return nil, lib.NewError(fmt.Errorf("content type %s not supported", t), lib.UnsupportedContentType)
	}

	return data, nil
}
