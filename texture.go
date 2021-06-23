package go3mx

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

const DefaultQuality = 80

func stringIn(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func isJpeg(h string) bool {
	return h[6:10] == "JFIF"
}

func isPng(h string) bool {
	return h[:8] == "\211PNG\r\n\032\n"
}

func peekImageFormat(buf string) string {
	if isJpeg(buf) {
		return "jpeg"
	} else if isPng(buf) {
		return "png"
	}
	return ""
}

type Texture struct {
	Image  image.Image
	ID     string
	Format string
}

func NewTexture(img image.Image, id string, format string) *Texture {
	return &Texture{Image: img, ID: id, Format: format}
}

func (t *Texture) Resource() Resource {
	return Resource{Type: TextureBufferType, Id: t.ID, Format: t.Format}
}

func (t *Texture) GetID() string {
	return t.ID
}

func (t *Texture) Marshal() []byte {
	writer := &bytes.Buffer{}
	if t.Format == "" {
		t.Format = "jpg"
	}
	encodeImage(t.Format, writer, t.Image)
	return writer.Bytes()
}

func (t *Texture) Unmarshal(buf []byte) {
	t.Format = peekImageFormat(string(buf[:10]))
	reader := bytes.NewBuffer(buf)
	t.Image = decodeImage(t.Format, reader)
}

func encodeImage(inputName string, writer io.Writer, rgba image.Image) {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		jpeg.Encode(writer, rgba, &jpeg.Options{Quality: DefaultQuality})
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(writer, rgba)
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
	}
	return nil
}
