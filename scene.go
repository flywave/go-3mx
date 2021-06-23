package go3mx

import (
	"encoding/json"
	"io"
)

type Scene struct {
	Version     int               `json:"3mxVersion"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Logo        string            `json:"logo"`
	Options     map[string]string `json:"sceneOptions"`
	Layers      []Layer           `json:"layers"`
}

type Layer struct {
	Type        string     `json:"type"`
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SRS         string     `json:"SRS"`
	SRSOrigin   [3]float64 `json:"SRSOrigin"`
	Root        string     `json:"root"`
}

func (o *Scene) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func SceneFromJson(data io.Reader) *Scene {
	var o *Scene
	json.NewDecoder(data).Decode(&o)
	return o
}
