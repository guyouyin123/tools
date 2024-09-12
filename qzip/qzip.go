package qzip

import (
	"archive/zip"
	"fmt"
	"github.com/guyouyin123/tools/qstring"
	"io"
	"os"
)

// AddFilesToZip 批量将多个文件添加到 ZIP 文件
func AddFilesToZip(filePaths []string, zipFileName, dirName string) error {
	// 创建 ZIP 文件
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		fmt.Println("创建 ZIP 文件失败:", err)
		return err
	}
	defer zipFile.Close()
	// 创建 ZIP Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		l := qstring.RsplitN(filePath, "/", 1)
		fileName := l[0]
		if len(l) == 2 {
			fileName = l[1]
		}
		customName := fmt.Sprintf("%s/%s", dirName, fileName)
		// 打开要压缩的文件
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		// 创建 ZIP 文件头
		zipFileWriter, err := zipWriter.Create(customName)
		if err != nil {
			return err
		}
		// 将文件内容复制到 ZIP 文件中
		_, err = io.Copy(zipFileWriter, file)
		if err != nil {
			return err
		}
	}
	return nil
}
