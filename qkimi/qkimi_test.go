package qkimi

import (
	"fmt"
	"strings"
	"testing"
)

var (
	apiKey      = "sk-rP7pbTplMHV6A0T8cGLLH2xH9eUD9NI6InwBFqSdijatmrhO"
	model       = "moonshot-v1-8k"
	temperature = float32(0.1)
)

func Test_ImagePath(t *testing.T) {
	//从图片中提取信息
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)

	filePath := "/图片9.png"
	isDelete := false //测试可以为false，产线必须为true

	cInfo, err := kimi.CreateFilePath(filePath)
	if err != nil {
		return
	}
	fileContent, err := kimi.GetFileContent(cInfo.Id)
	if err != nil {
		return
	}
	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "你是一个非常有经验的职员，你可以根据自己的经验从文本中提取工号，工号的特征为字母+数组组合，在4～10位数以内",
		},
		{
			"role":    "system",
			"content": fileContent.Content,
		},
		{"role": "user", "content": "帮我从文字中提取工号,以workNumber=形式输出"},
	}
	chatData, err := kimi.Chat(messages)
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

func Test_ImageUrl(t *testing.T) {
	imgUrl := "https://img2024.cnblogs.com/blog.jpg"
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)

	cInfo, err := kimi.CreateFileUrl(imgUrl)
	if err != nil {
		return
	}
	fileContent, err := kimi.GetFileContent(cInfo.Id)
	if err != nil {
		return
	}
	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "你是一个非常有经验的职员，你可以根据自己的经验从文本中提取工号，工号的特征为字母+数组组合，在4～10位数以内",
		},
		{
			"role":    "system",
			"content": fileContent.Content,
		},
		{"role": "user", "content": "帮我从文字中提取工号,以workNumber=形式输出"},
	}
	chatData, err := kimi.Chat(messages)
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
	if false {
		//删除文件，每次请求将消耗每分钟的RPM。每个账号最多上传1000个文件。所以必须删除操作
		_, _ = kimi.DeleteFile(cInfo.Id)
	}
	fmt.Println(workNumber)
}

func Test_UsersMeBalance(t *testing.T) {
	//查询余额
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)
	resp, err := kimi.UsersMeBalance()
	if err != nil {
		return
	}
	fmt.Println(resp)
}

func Test_chat(t *testing.T) {
	//发起会话
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)

	messages := []map[string]interface{}{
		{
			"role":    "system",
			"content": "你是一名经验丰富的金融资本家，你具有多年的炒股经验",
		},
		{"role": "user", "content": "告诉我如何选股"},
	}
	resp, err := kimi.Chat(messages)
	if err != nil {
		return
	}
	fmt.Println(resp)
}

func Test_DeleteFiles(t *testing.T) {
	//清空文件
	kimi := &KiMiAi{}
	kimi.InitKiMiAi(apiKey, model, temperature)
	kimi.DeleteFiles()
}
