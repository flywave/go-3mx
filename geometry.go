package go3mx

import (
	_ "github.com/flywave/go-ctm"
)

type Geometry struct{}

func (g *Geometry) Marshalling() []byte {
	return nil
}
