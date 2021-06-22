package go3mx

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
