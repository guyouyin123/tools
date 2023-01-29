package qtable

import (
	"bytes"
	"io"
	"github.com/guyouyin123/tools/qcsv"
)

type csvRead struct {
}

func newCsvRead() *csvRead {
	return &csvRead{}
}

func (c csvRead) FileType() string {
	return "csv"
}

func (c csvRead) Read(reader io.Reader) (data []map[string]string, err error) {
	var buf bytes.Buffer
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return
	}
	return qcsv.Read(buf.Bytes())
}
