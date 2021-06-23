package go3mx

import (
	"encoding/json"
	"io"
)

const (
	TextureBufferType  = "textureBuffer"
	GeometryBufferType = "geometryBuffer"
	TextureFileType    = "textureFile"
	GeometryFileType   = "geometryFile"
)

type Node struct {
	Id                string     `json:"id"`
	BBoXMin           [3]float64 `json:"bbMin"`
	BBoXMax           [3]float64 `json:"bbMax"`
	MaxScreenDiameter float64    `json:"maxScreenDiameter"`
	Children          []string   `json:"children"`
	Resources         []string   `json:"resources"`
}

type Resource struct {
	Type      string      `json:"type"`
	Id        string      `json:"id"`
	Format    string      `json:"format"`
	Size      uint32      `json:"size"`
	BBoXMin   *[3]float64 `json:"bbMin,omitempty"`
	BBoXMax   *[3]float64 `json:"bbMax,omitempty"`
	Texture   *string     `json:"texture,omitempty"`
	PointSize *uint32     `json:"pointSize,omitempty"`
	File      *string     `json:"file,omitempty"`
}

type Header struct {
	Version   uint32     `json:"version"`
	Nodes     []Node     `json:"nodes"`
	Resources []Resource `json:"resources"`
}

func (o *Header) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func HeaderFromJson(data io.Reader) *Header {
	var o *Header
	json.NewDecoder(data).Decode(&o)
	return o
}
