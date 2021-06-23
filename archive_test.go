package go3mx

import (
	"os"
	"testing"
)

func Test3mxbReader(t *testing.T) {
	a := &Archive{}
	f, _ := os.Open("./testdata/Model.3mxb")

	a.Load(f)

	f.Close()

}
