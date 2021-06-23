package go3mx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const (
	MAGIC_NUMBER = "3MXBO"
)

var byteorder = binary.LittleEndian

type Archive struct {
	Magic    []byte
	Size     int32
	Header   *Header
	geos     []Geometry
	texs     []Texture
	buffers  [][]byte
	basePath string
}

func NewArchive(nodes []Node, geos []Geometry, texs []Texture) *Archive {
	a := &Archive{Magic: []byte(MAGIC_NUMBER)}
	a.buffers = a.pack(nodes, geos, texs)
	return a
}

func (a *Archive) pack(nodes []Node, geos []Geometry, texs []Texture) [][]byte {
	a.Header = new(Header)
	a.Header.Nodes = nodes
	a.Header.Resources = make([]Resource, 0, len(geos)+len(texs))
	buffers := make([][]byte, 0, len(geos)+len(texs))
	off := 0
	for i := range texs {
		res := texs[i].Resource()
		buffers = append(buffers, texs[i].Marshal())
		res.Size = uint32(len(buffers[len(buffers)-1]))
		a.Header.Resources = append(a.Header.Resources, res)
		off++
	}

	for i := range geos {
		res := geos[i].Resource()
		buffers = append(buffers, geos[i].Marshal())
		res.Size = uint32(len(buffers[len(buffers)-1]))
		a.Header.Resources = append(a.Header.Resources, res)
		off++
	}
	return buffers
}

func (a *Archive) packHeader() ([]byte, uint32) {
	data := []byte(a.Header.ToJson())
	return data, uint32(len(data))
}

func (a *Archive) unpackResource(res Resource, rd io.Reader) error {
	if res.Type == TextureBufferType {
		buf := make([]byte, res.Size)
		_, err := rd.Read(buf)
		if err != nil {
			return err
		}
		tex := &Texture{}
		tex.Unmarshal(buf)
		a.texs = append(a.texs, *tex)
	} else if res.Type == GeometryBufferType {
		buf := make([]byte, res.Size)
		_, err := rd.Read(buf)
		if err != nil {
			return err
		}
		if res.Format == FORMAT_CTM {
			geom := &Mesh{}
			geom.Unmarshal(buf)
			a.geos = append(a.geos, geom)
		} else if res.Format == FORMAT_XYZ {
			geom := &PointCloud{}
			geom.Unmarshal(buf)
			a.geos = append(a.geos, geom)
		}
	} else if res.Type == TextureFileType {
		path := path.Join(a.basePath, *res.File)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		tex := &Texture{}
		tex.Unmarshal(buf)
		a.texs = append(a.texs, *tex)
	} else if res.Type == GeometryFileType {
		path := path.Join(a.basePath, *res.File)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if res.Format == FORMAT_CTM {
			geom := &Mesh{}
			geom.Unmarshal(buf)
			a.geos = append(a.geos, geom)
		} else if res.Format == FORMAT_XYZ {
			geom := &PointCloud{}
			geom.Unmarshal(buf)
			a.geos = append(a.geos, geom)
		}
	}
	return nil
}

func (a *Archive) Save(wr io.Writer) error {
	if a.buffers == nil {
		return errors.New("empty archive")
	}
	_, err := wr.Write(a.Magic)
	if err != nil {
		return err
	}
	header, hsi := a.packHeader()
	err = binary.Write(wr, byteorder, hsi)
	if err != nil {
		return err
	}
	_, err = wr.Write(header)
	if err != nil {
		return err
	}
	for _, bf := range a.buffers {
		_, err = wr.Write(bf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Archive) Load(rd io.Reader) error {
	magic := make([]byte, 5)
	_, err := rd.Read(magic)
	if err != nil {
		return err
	}
	if string(magic) != MAGIC_NUMBER {
		return errors.New("file format is not 3dmxb")
	} else {
		a.Magic = magic
	}

	err = binary.Read(rd, byteorder, &a.Size)
	if err != nil {
		return err
	}

	hbuf := make([]byte, a.Size)
	_, err = rd.Read(hbuf)
	if err != nil {
		return err
	}

	a.Header = HeaderFromJson(bytes.NewBuffer(hbuf))

	for _, res := range a.Header.Resources {
		err = a.unpackResource(res, rd)
		if err != nil {
			return err
		}
	}

	return nil
}
