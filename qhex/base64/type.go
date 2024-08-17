package base64

type PageToken struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	Timestamp int64
}
