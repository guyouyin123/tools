package qkimi

import (
	"bytes"
	"encoding/json"
	"fmt"
	qhttp "github.com/guyouyin123/tools/qhttp"
	jsoniter "github.com/json-iterator/go"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (kimi *KiMiAi) InitKiMiAi(apiKey, model string, temperature float32) {
	kimi.BaseUrl = "https://api.moonshot.cn/v1"
	kimi.Header = map[string]interface{}{
		"Authorization": "Bearer " + apiKey,
	}
	kimi.Model = model
	kimi.Temperature = temperature
}

/*
Chat 会话
messages:消息体
model:会话模型 目前支持:moonshot-v1-8k/1M/12¥ moonshot-v1-32k/1M/24¥ moonshot-v1-128k/1M/60¥
temperature:使用什么采样温度，介于 0 和 1 之间。较高的值（如 0.7）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定性。
*/
func (kimi *KiMiAi) Chat(messages []map[string]interface{}) (*ChatStruct, error) {
	url := kimi.BaseUrl + "/chat/completions"
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":       kimi.Model,
		"messages":    messages,
		"temperature": kimi.Temperature,
	})
	header := kimi.Header
	header["Content-Type"] = "application/json"
	resp, err := qhttp.Post(url, nil, kimi.Header, nil, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	data := &ChatStruct{}
	_ = jsoniter.Unmarshal(resp, &data)
	return data, nil
}

// CreateFilePath 上传本地文件
func (kimi *KiMiAi) CreateFilePath(filePath string) (*FileCreate, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	formField, err := writer.CreateFormField("purpose")
	if err != nil {
		return nil, err
	}
	_, err = formField.Write([]byte("file-extract"))

	fileName := filepath.Base(filePath)
	if getImageExtension(filePath) == "" {
		fileName += ".jpg"
	}
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	fd, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
		return nil, err
	}
	writer.Close()

	url := kimi.BaseUrl + "/files"
	header := kimi.Header
	header["Content-Type"] = writer.FormDataContentType()
	resp, err := qhttp.Post(url, nil, kimi.Header, nil, form)
	if err != nil {
		return nil, err
	}
	data := &FileCreate{}
	_ = jsoniter.Unmarshal(resp, &data)
	return data, nil
}

// CreateFileUrl 上传网络文件
func (kimi *KiMiAi) CreateFileUrl(fileUrl string) (*FileCreate, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)

	// 添加表单字段
	formField, err := writer.CreateFormField("purpose")
	if err != nil {
		return nil, err
	}
	_, err = formField.Write([]byte("file-extract"))
	if err != nil {
		return nil, err
	}

	// 创建文件字段
	fileName := filepath.Base(fileUrl)
	if getImageExtension(fileUrl) == "" {
		fileName += ".jpg"
	}
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	// 发送 GET 请求以获取文件
	response, err := http.Get(fileUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching file: %w", err)
	}
	defer response.Body.Close()

	// 检查响应状态
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", response.StatusCode)
	}

	// 将响应体写入文件字段
	_, err = io.Copy(fw, response.Body)
	if err != nil {
		return nil, err
	}

	// 关闭 writer
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// 发送 POST 请求
	url := kimi.BaseUrl + "/files"
	header := kimi.Header
	header["Content-Type"] = writer.FormDataContentType()

	resp, err := qhttp.Post(url, nil, header, nil, form)
	if err != nil {
		return nil, err
	}

	data := &FileCreate{}
	if err := jsoniter.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// deleteFile 删除文件
func (kimi *KiMiAi) DeleteFile(fileId string) (*DeleteFile, error) {
	url := kimi.BaseUrl + "/files/" + fileId
	resp, err := qhttp.Delete(url, nil, kimi.Header, nil, nil)
	if err != nil {
		return nil, err
	}
	data := &DeleteFile{}
	_ = jsoniter.Unmarshal(resp, &data)
	return data, nil
}

// GetFileContent 图片提取文字
func (kimi *KiMiAi) GetFileContent(fileID string) (*FileContent, error) {
	url := kimi.BaseUrl + "/files/" + fileID + "/content"
	resp, err := qhttp.Get(url, nil, kimi.Header)
	if err != nil {
		return nil, err
	}
	data := &FileContent{}
	_ = jsoniter.Unmarshal(resp, &data)
	return data, nil
}

// UsersMeBalance 查询余额
func (kimi *KiMiAi) UsersMeBalance() (*Balance, error) {
	url := kimi.BaseUrl + "/users/me/balance"
	resp, err := qhttp.Get(url, nil, kimi.Header)
	if err != nil {
		return nil, err
	}
	data := &Balance{}
	_ = jsoniter.Unmarshal(resp, &data)
	return data, nil
}

func getImageExtension(url string) string {
	ext := path.Ext(url)
	if ext == "" {
		if idx := strings.Index(url, "?"); idx != -1 {
			url = url[:idx]
		}
		ext = path.Ext(url)
	}
	return ext
}
