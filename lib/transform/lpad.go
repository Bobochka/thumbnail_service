package transform

import (
	"crypto/sha1"
	"fmt"
	"image"

	"image/draw"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/nfnt/resize"
)

type LPad struct {
	Width  int
	Height int
	codec  Img
}

func NewLPad(width, height int) *LPad {
	return &LPad{
		codec:  Img{},
		Width:  width,
		Height: height,
	}
}

func (t LPad) Fingerprint(data []byte) string {
	return fmt.Sprintf("%x_%v_%v", sha1.Sum(data), t.Width, t.Height)
}

func (t LPad) Perform(data []byte) ([]byte, error) {
	img, err := t.codec.Decode(data)
	if err != nil {
		return nil, lib.NewError(err, lib.UnsupportedContentType)
	}

	img, err = t.perform(img)
	if err != nil {
		return nil, lib.NewError(err, lib.TransformationFailure)
	}

	imgBytes, err := t.codec.Encode(img)
	if err != nil {
		return nil, lib.NewError(err, lib.EncodingFailure)
	}

	return imgBytes, nil
}

func (t *LPad) perform(img image.Image) (image.Image, error) {
	rect := img.Bounds()
	origW := rect.Dx()
	origH := rect.Dy()

	if origH == t.Height && origW == t.Width {
		return img, nil
	}

	thumb := resize.Thumbnail(uint(t.Width), uint(t.Height), img, 0)

	if t.isScaledDownsize(origW, origH) {
		return thumb, nil
	}

	pt := t.findSp(thumb.Bounds())

	dst := image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))
	draw.Draw(dst, dst.Bounds(), thumb, pt, draw.Src)

	return dst, nil
}

func (t LPad) isScaledDownsize(origW, origH int) bool {
	if t.Width >= origW || t.Height >= origH {
		return false
	}

	newHeight := float64(origH) * float64(t.Width) / float64(origW)

	roundingDelta := 0.5
	return int(newHeight+roundingDelta) == t.Height
}

func (t LPad) findSp(thumbRect image.Rectangle) image.Point {
	var ptX, ptY int

	ptX = max((t.Width-thumbRect.Dx())/2, 0)
	ptY = max((t.Height-thumbRect.Dy())/2, 0)

	return image.Pt(-ptX, -ptY)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
