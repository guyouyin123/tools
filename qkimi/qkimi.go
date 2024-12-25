package qkimi

import (
	"bytes"
	"encoding/json"
	qhttp "github.com/guyouyin123/tools/qhttp"
	jsoniter "github.com/json-iterator/go"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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

// createFile 上传文件
func (kimi *KiMiAi) CreateFile(filePath string) (*FileCreate, error) {
	form := new(bytes.Buffer)
	writer := multipart.NewWriter(form)
	formField, err := writer.CreateFormField("purpose")
	if err != nil {
		return nil, err
	}
	_, err = formField.Write([]byte("file-extract"))
	fw, err := writer.CreateFormFile("file", filepath.Base(filePath))
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
