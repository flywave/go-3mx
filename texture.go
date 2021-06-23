package go3mx

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type Texture struct{}

func (g *Texture) Marshalling() []byte {
	return nil
}

func encodeImage(inputName string, writer io.Writer, rgba image.Image) {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		jpeg.Encode(writer, rgba, nil)
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(writer, rgba)
	} else if strings.HasSuffix(inputName, "gif") {
		gif.Encode(writer, rgba, nil)
	}
}

func decodeImage(inputName string, reader io.Reader) image.Image {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		img, err := jpeg.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	} else if strings.HasSuffix(inputName, "png") {
		img, err := png.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	} else if strings.HasSuffix(inputName, "gif") {
		img, err := gif.Decode(reader)
		if err != nil {
			return nil
		}
		return img
	}
	return nil
}
