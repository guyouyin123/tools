[toc]

# kimi-AI

## 说明

```go
本文主要通过kimi提供的api实现golang版本sdk，使用时需要注意速率限制，根据自己的版本查看对应速率
```

## 开发文档

```go
https://platform.moonshot.cn/docs/intro
```

## 实现的api

```go
本文实现的api：
上传文件：https://api.moonshot.cn/v1/files
查询所有文件：https://api.moonshot.cn/v1/files
删除文件：https://api.moonshot.cn/v1/files/{fileId}
删除所有文件：DeleteFiles
图文识别：https://api.moonshot.cn/v1/files/{fileId}/content
ai对话：https://api.moonshot.cn/v1/chat/completions
查询余额：https://api.moonshot.cn/v1/users/me/balance
```

## test-从图片中提取信息

```go
//本地文件
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

//网络文件
func Test_ImageUrl(t *testing.T) {
	imgUrl := "https://img2024.cnblogs.com/blog/1736414/202412/1736414-20241227101453173-390339808.png"
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
```

## test-发起会话

```go
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
```

## 清空文件
```go
func Test_DeleteFiles(t *testing.T) {
    //清空文件
    kimi := &KiMiAi{}
    kimi.InitKiMiAi(apiKey, model, temperature)
    kimi.DeleteFiles()
}

```

