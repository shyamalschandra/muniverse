package muniverse

import (
	"bytes"
	"image"
	"image/png"

	"github.com/unixpickle/essentials"
)

// Obs is an observation from an environment.
type Obs interface {
	Image() (image.Image, error)
}

// RGB generates an RGB pixel buffer for an observation.
//
// The buffer is packed the same way as the pixel data in
// *image.RGBA, but without an alpha channel.
func RGB(o Obs) (buffer []uint8, width, height int, err error) {
	img, err := o.Image()
	if err != nil {
		return
	}
	width, height = img.Bounds().Dx(), img.Bounds().Dy()
	switch img := img.(type) {
	case *image.NRGBA:
		if img.Stride != width*4 {
			buffer = naiveRGB(img)
		} else {
			buffer = rgbaToRGB(img.Pix)
		}
	case *image.RGBA:
		if img.Stride != width*4 {
			buffer = naiveRGB(img)
		} else {
			buffer = rgbaToRGB(img.Pix)
		}
	default:
		buffer = naiveRGB(img)
	}
	return
}

func rgbaToRGB(rgba []uint8) []uint8 {
	buffer := make([]uint8, 3*(len(rgba)/4))
	var destIdx int
	for i := 0; i < len(rgba); i += 4 {
		buffer[destIdx] = rgba[i]
		buffer[destIdx+1] = rgba[i+1]
		buffer[destIdx+2] = rgba[i+2]
		destIdx += 3
	}
	return buffer
}

func naiveRGB(img image.Image) []uint8 {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	min := img.Bounds().Min
	buffer := make([]uint8, width*height*3)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x+min.X, y+min.Y).RGBA()
			idx := (x + y*width) * 3
			buffer[idx] = uint8(r >> 8)
			buffer[idx+1] = uint8(g >> 8)
			buffer[idx+2] = uint8(b >> 8)
		}
	}
	return buffer
}

type pngObs []byte

func (p pngObs) Image() (img image.Image, err error) {
	defer essentials.AddCtxTo("decode PNG observation", &err)
	return png.Decode(bytes.NewReader(p))
}