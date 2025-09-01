package qkimi

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

var (
	apiKey = ""
	//model       = "moonshot-v1-8k"
	model       = "moonshot-v1-8k-vision-preview"
	temperature = float32(0.1)
)

func Test_ImagePath(t *testing.T) {
	//从图片中提取信息
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)

	filePath := "/图片9.png"
	isDelete := false //测试可以为false，产线必须为true

	cInfo, err := kimi.CreateFileUrl(filePath, 5)
	if err != nil {
		return
	}
	fileContent, err := kimi.GetFileContent(cInfo.Id, 5)
	if err != nil {
		return
	}

	imageContent := fileContent.Content

	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "你是一个非常有经验的职员，你可以根据自己的经验从文本中提取工号，工号的特征为字母+数组组合，在4～10位数以内",
		},
		{
			"role":    "system",
			"content": imageContent,
		},
		{"role": "user", "content": "帮我从文字中提取工号,以workNumber=形式输出"},
	}
	chatData, err := kimi.Chat(messages, 5)
	if err != nil {
		return
	}

	if len(chatData.Choices) == 0 {
		return
	}
	content := chatData.Choices[0].Message.Content
	contentSplit := strings.SplitN(content, "=", 2)
	if len(contentSplit) < 2 {
		return
	}
	workNumber := contentSplit[1]
	if isDelete {
		//删除文件，每次请求将消耗每分钟的RPM。每个账号最多上传1000个文件。所以必须删除操作
		_, _ = kimi.DeleteFile(cInfo.Id)
	}
	fmt.Println(workNumber)
}

func imageToBase64(url string) (string, error) {
	// 下载图片
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应体
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// 编码为Base64
	base64Str := base64.StdEncoding.EncodeToString(imageData)
	imageBase64Str := fmt.Sprintf("data:image/jpeg;base64,%s", base64Str)
	return imageBase64Str, nil
}

func Test_ImageUrl(t *testing.T) {
	imgUrl := "http://woda-app-public.oss-cn-shanghai.aliyuncs.com/zxx/WorkCard/288bd33f-cf57-4fe4-8be1-bff9f2dd1052" //ok
	//imgUrl := "http://woda-app-public.oss-cn-shanghai.aliyuncs.com/zxx/WorkCard/5681028c-b0e2-4397-be1c-18ce13ede8a5" //no
	//imgUrl := "http://woda-app-public.oss-cn-shanghai.aliyuncs.com/zxx/WorkCard/C93D1621-369D-45C2-AA52-1A9E03DAFF92.jpg" //no
	//imgUrl := "http://woda-app-public.oss-cn-shanghai.aliyuncs.com/zxx/WorkCard/2b9996a9-ab85-41b2-bdc5-e071a8c9c9f6" //ok
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)
	base64Image, _ := imageToBase64(imgUrl)
	messages := []map[string]interface{}{
		{
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "image_url",
					"image_url": map[string]interface{}{
						"url": base64Image,
					},
				},
				{
					"type": "text",
					"text": "读取图片内容。",
				},
			},
		},
		{"role": "user", "content": "帮我从文字中提取工号,以workNumber=形式输出"},
	}

	chatData, err := kimi.Chat(messages, 5)
	if err != nil {
		return
	}
	if len(chatData.Choices) == 0 {
		return
	}
	content := chatData.Choices[0].Message.Content
	contentSplit := strings.SplitN(content, "=", 2)
	if len(contentSplit) < 2 {
		return
	}
	workNumber := contentSplit[1]
	fmt.Println(workNumber)
}
