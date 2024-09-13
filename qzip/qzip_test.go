package qzip

import (
	"fmt"
	"testing"
)

func TestAddFileToZip(t *testing.T) {
	files := []string{"/Users/jeff/myself/tools/qzip/test1你好.txt", "/Users/jeff/myself/tools/qzip/test2.txt"}
	zipFileName := "/Users/jeff/myself/tools/qzip/导出数据.zip" //压缩包名称
	dirName := "123"                                        //压缩包内部目录，内部无需分目录为空字符

	err := AddFilesToZip(files, zipFileName, dirName)
	if err != nil {
		return
	}

	fmt.Println("压缩完成:", zipFileName)
}

func Test_Unzip(t *testing.T) {
	zipFile := "/Users/jeff/myself/tools/qzip/导出数据.zip" //压缩包名称
	destDir := "./output"                               // 解压缩目标目录

	err := Unzip(zipFile, destDir)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Unzip completed successfully.")
	}
}
