package transform

import (
	"bytes"
	"errors"
	"image"
	_ "image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
)

var ErrUnknownFormat = errors.New("can't decode: unknown image format")

type Img struct {
}

func (Img) Encode(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (Img) Decode(data []byte) (image.Image, error) {
	r := bytes.NewReader(data)

	img, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	if format == "" || img == nil {
		return nil, ErrUnknownFormat
	}

	return img, nil
}
