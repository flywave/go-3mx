package go3mx

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"

	"github.com/flywave/go-ctm"

	"github.com/flywave/go3d/vec2"
	"github.com/flywave/go3d/vec3"
)

const (
	FORMAT_CTM = "ctm"
	FORMAT_XYZ = "xyz"
)

type Color [4]uint8

type Geometry interface {
	Marshal() []byte
	Unmarshal(buf []byte)
	BoundBox() vec3.Box
	GetID() string
	Format() string
	Resource() Resource
}

type Mesh struct {
	Geometry
	Vertices []vec3.T
	Normals  []vec3.T
	UVCoords []vec2.T
	Indices  [][3]uint32
	ID       string
	Texture  string
}

func NewMesh(id string, vertices []vec3.T, indices [][3]uint32, normals []vec3.T, uvs []vec2.T, texture string) Geometry {
	return &Mesh{Vertices: vertices, Normals: normals, UVCoords: uvs, Indices: indices, ID: id, Texture: texture}
}

func (m *Mesh) toCTM() *ctm.Mesh {
	cm := ctm.NewEmptyMesh()
	cm.AddUVMap(m.UVCoords, "", m.Texture)
	return cm
}

func (m *Mesh) Resource() Resource {
	bb := m.BoundBox()
	min := [3]float32{bb.Min[0], bb.Min[1], bb.Min[2]}
	max := [3]float32{bb.Max[0], bb.Max[1], bb.Max[2]}

	return Resource{Type: GeometryBufferType, Id: m.ID, Format: FORMAT_CTM, BBoXMin: &min, BBoXMax: &max, Texture: &m.Texture}
}

func (m *Mesh) Format() string {
	return FORMAT_CTM
}

func (m *Mesh) Marshal() []byte {
	cm := m.toCTM()
	return cm.GetContext().SaveToBuffer()
}

func (m *Mesh) Unmarshal(buf []byte) {
	cm := ctm.NewEmptyMesh()
	cm.GetContext().LoadFromBuffer(buf)

	vers := cm.GetVertices()
	m.Vertices = make([]vec3.T, len(vers))

	copy(m.Vertices, vers)
	if cm.HasNormals() {
		m.Normals = make([]vec3.T, len(vers))
		copy(m.Normals, cm.GetNormals())
	}
	faces := cm.GetFaces()
	m.Indices = make([][3]uint32, len(faces))
	copy(m.Indices, faces)

	if cm.GetUVMapCount() > 0 {
		size := cm.GetVertCount()
		data := cm.GetContext().GetFloatArray(ctm.CTM_UV_MAP_1)

		var bufSlice []vec2.T
		bufHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bufSlice)))
		bufHeader.Cap = int(size)
		bufHeader.Len = int(size)
		bufHeader.Data = uintptr(unsafe.Pointer(data))

		m.UVCoords = make([]vec2.T, size)

		copy(m.UVCoords, bufSlice)
	}
}

func (m *Mesh) BoundBox() vec3.Box {
	ret := vec3.Box{Min: vec3.MaxVal, Max: vec3.MinVal}
	for _, v := range m.Vertices {
		ret.Extend(&v)
	}
	return vec3.Box{}
}

func (m *Mesh) GetID() string {
	return m.ID
}

type PointCloud struct {
	Geometry
	ID        string
	PointSize uint32
	Vertices  []vec3.T
	Colors    []Color
}

func NewPointCloud(id string, vertices []vec3.T, colors []Color, pointSize uint32) Geometry {
	return &PointCloud{Vertices: vertices, Colors: colors, PointSize: pointSize, ID: id}
}

func (m *PointCloud) Resource() Resource {
	bb := m.BoundBox()
	min := [3]float32{bb.Min[0], bb.Min[1], bb.Min[2]}
	max := [3]float32{bb.Max[0], bb.Max[1], bb.Max[2]}

	return Resource{Type: GeometryBufferType, Id: m.ID, Format: FORMAT_XYZ, BBoXMin: &min, BBoXMax: &max, PointSize: &m.PointSize}
}

func (m *PointCloud) Format() string {
	return FORMAT_XYZ
}

func (m *PointCloud) Marshal() []byte {
	writer := &bytes.Buffer{}

	si := int32(len(m.Vertices))

	var vSlice []float32
	vHeader := (*reflect.SliceHeader)((unsafe.Pointer(&vSlice)))
	vHeader.Cap = int(si * 3)
	vHeader.Len = int(si * 3)
	vHeader.Data = uintptr(unsafe.Pointer(&m.Vertices))

	var cSlice []uint8
	cHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cSlice)))
	cHeader.Cap = int(si * 3)
	cHeader.Len = int(si * 3)
	cHeader.Data = uintptr(unsafe.Pointer(&m.Colors))

	err := binary.Write(writer, byteorder, si)
	if err != nil {
		return nil
	}
	err = binary.Write(writer, byteorder, vSlice)
	if err != nil {
		return nil
	}
	err = binary.Write(writer, byteorder, cSlice)
	if err != nil {
		return nil
	}
	return writer.Bytes()
}

func (m *PointCloud) Unmarshal(buf []byte) {
	reader := bytes.NewBuffer(buf)
	var si int32
	err := binary.Read(reader, byteorder, si)
	if err != nil {
		return
	}
	m.Vertices = make([]vec3.T, si)

	var vSlice []float32
	vHeader := (*reflect.SliceHeader)((unsafe.Pointer(&vSlice)))
	vHeader.Cap = int(si * 3)
	vHeader.Len = int(si * 3)
	vHeader.Data = uintptr(unsafe.Pointer(&m.Vertices))

	err = binary.Read(reader, byteorder, vSlice)
	if err != nil {
		return
	}

	m.Colors = make([]Color, si)

	var cSlice []uint8
	cHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cSlice)))
	cHeader.Cap = int(si * 3)
	cHeader.Len = int(si * 3)
	cHeader.Data = uintptr(unsafe.Pointer(&m.Colors))

	err = binary.Read(reader, byteorder, cSlice)
	if err != nil {
		return
	}
}

func (m *PointCloud) BoundBox() vec3.Box {
	ret := vec3.Box{Min: vec3.MaxVal, Max: vec3.MinVal}
	for _, v := range m.Vertices {
		ret.Extend(&v)
	}
	return vec3.Box{}
}

func (m *PointCloud) GetID() string {
	return m.ID
}
