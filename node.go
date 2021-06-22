package go3mx

const (
	MAGIC_NUMBER = "3MXBO"
)

type File struct {
	Magic []byte
	Size  uint32
}

type Node struct {
	Id                string     `json:"id"`
	BBoXMin           [3]float64 `json:"bbMin"`
	BBoXMax           [3]float64 `json:"bbMax"`
	MaxScreenDiameter float64    `json:"maxScreenDiameter"`
	Children          []string   `json:"children"`
	Resources         []string   `json:"resources"`
}

const (
	TextureBufferType  = "textureBuffer"
	GeometryBufferType = "geometryBuffer"
	TextureFileType    = "textureFile"
	GeometryFileType   = "geometryFile"
)

type Resource struct {
	Type    string      `json:"type"`
	Id      string      `json:"id"`
	Format  string      `json:"format"`
	Size    uint32      `json:"size"`
	BBoXMin *[3]float64 `json:"bbMin,omitempty"`
	BBoXMax *[3]float64 `json:"bbMax,omitempty"`
	Texture *string     `json:"texture,omitempty"`
	File    *string     `json:"file,omitempty"`
}

type Header struct {
	Version   uint32     `json:"version"`
	Nodes     []Node     `json:"nodes"`
	Resources []Resource `json:"resources"`
}
