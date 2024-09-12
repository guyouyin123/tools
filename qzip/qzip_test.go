package qzip

import (
	"fmt"
	"testing"
)

func TestAddFileToZip(t *testing.T) {
	files := []string{"/Users/jeff/myself/tools/qzip/test1你好.txt", "/Users/jeff/myself/tools/qzip/test2.txt"}
	zipFileName := "/Users/jeff/myself/tools/qzip/output.zip" //压缩包名称
	dirName := "Jeff"                                         //压缩包内部目录

	err := AddFilesToZip(files, zipFileName, dirName)
	if err != nil {
		return
	}

	fmt.Println("压缩完成:", zipFileName)
}
