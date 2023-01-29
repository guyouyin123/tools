package qtable

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

func TestTableReader(t *testing.T) {
	filename := "./template.xlsx"
	rf, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	table, err := TableRead(path.Ext(filename))
	if err != nil {
		t.Error(err)
		return
	}
	newRd := bytes.NewReader(rf)
	data, err := table.Read(newRd)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("data", data)
}
