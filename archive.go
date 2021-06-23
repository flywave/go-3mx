package go3mx

const (
	MAGIC_NUMBER = "3MXBO"
)

type Archive struct {
	Magic   []byte
	Size    uint32
	Header  *Header
	Buffers [][]byte
}

func NewArchive() *Archive {
	return &Archive{Magic: []byte(MAGIC_NUMBER)}
}

func (a *Archive) getHeader() ([]byte, uint32) {
	data := []byte(a.Header.ToJson())
	return data, uint32(len(data))
}

func (a *Archive) appendGeometry(g *Geometry) (int, uint32) {
	a.Buffers = append(a.Buffers, g.Marshalling())
	return len(a.Buffers) - 1, uint32(len(a.Buffers))
}

func (a *Archive) appendTexture(t *Texture) (int, uint32) {
	a.Buffers = append(a.Buffers, t.Marshalling())
	return len(a.Buffers) - 1, uint32(len(a.Buffers))
}
