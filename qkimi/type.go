package qkimi

type KiMiAi struct {
	BaseUrl     string                 //https://api.moonshot.cn/v1
	Header      map[string]interface{} //请求头
	Model       string                 //会话模型 目前支持:moonshot-v1-8k/1M/12¥ moonshot-v1-32k/1M/24¥ moonshot-v1-128k/1M/60¥
	Temperature float32                //使用什么采样温度，介于 0 和 1 之间。较高的值（如 0.7）将使输出更加随机，而较低的值（如 0.2）将使其更加集中和确定性。
}

type FileContent struct {
	Content  string `json:"content"`
	FileType string `json:"file_type"`
	FileName string `json:"file_name"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type FileCreate struct {
	Id            string `json:"id"`
	Object        string `json:"object"`
	Bytes         int64  `json:"bytes"`
	CreatedAt     int64  `json:"created_at"`
	Filename      string `json:"filename"`
	Purpose       string `json:"purpose"`
	Status        string `json:"status"`
	StatusDetails string `json:"status_details"`
}

type ChatStruct struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}
		FinishReason string `json:"finishReason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"` //消耗总token
	}
}

type DeleteFile struct {
	Deleted bool   `json:"deleted"`
	Id      string `json:"id"`
	Object  string `json:"object"`
}

type Balance struct {
	Code int `json:"code"`
	Data struct {
		AvailableBalance float64 `json:"available_balance"` //可用余额，包括现金余额和代金券余额, 当它小于等于 0 时, 用户不可调用推理 API
		VoucherBalance   float64 `json:"voucher_balance"`   //代金券余额, 不会为负数
		CashBalance      float64 `json:"cash_balance"`      //现金余额, 可能为负数, 代表用户欠费, 当它为负数时, available_balance 为 voucher_balance 的值
	}
	Scode  string `json:"scode"`
	Status string `json:"status"`
}
