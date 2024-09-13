package qzip

import (
	"archive/zip"
	"fmt"
	"github.com/guyouyin123/tools/qstring"
	"io"
	"os"
	"path/filepath"
)

// AddFilesToZip 批量将多个文件添加到 ZIP 文件
func AddFilesToZip(filePaths []string, zipFileName, dirName string) error {
	l := qstring.RsplitN(zipFileName, "/", 1)
	baseName := l[0]
	if len(l) == 2 {
		baseName = l[1]
	}

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
		customName := fmt.Sprintf("%s/%s/%s", baseName, dirName, fileName)
		if dirName == "" {
			customName = fmt.Sprintf("%s/%s", baseName, fileName)
		}
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

func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	// 创建目标目录
	os.MkdirAll(dest, 0755)
	// 遍历 ZIP 文件中的每个文件
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
		outFile.Close()
		rc.Close()
	}
	return nil
}
