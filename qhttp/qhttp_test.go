package qhttp

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestGet(t *testing.T) {
	urlR := "https://www.baidu.com/"
	header := map[string]interface{}{
		"content-type": "application/json",
		"access-token": "xxoo",
	}
	params := map[string]interface{}{
		"date_type": 30,
	}
	resp, _ := Get(urlR, params, header)
	fmt.Println(string(resp))
}

func TestPost(t *testing.T) {
	urlR := "https://open.douyin.com/api/douyin/v1/video/video_data/"
	header := map[string]interface{}{
		"content-type": "application/json",
		"access-token": "xx",
	}
	params := map[string]interface{}{
		"open_id": 123,
	}
	data := map[string]interface{}{
		"item_ids": 123,
	}
	respByte, _ := Post(urlR, params, header, data)
	fmt.Println(string(respByte))
	dto := DouYinVideoDataDto{}
	_ = jsoniter.Unmarshal(respByte, &dto)
	fmt.Println(dto)
}

type DouYinVideoDataDto struct {
	Data struct {
		List []struct {
			Cover       string `json:"cover"`        //视频封面
			CreateTime  int64  `json:"create_time"`  //视频创建时间戳
			IsReviewed  bool   `json:"total_like"`   //表示是否审核结束。审核通过或者失败都会返回true，审核中返回false。
			IsTop       bool   `json:"total_play"`   //是否置顶
			ItemId      string `json:"item_id"`      //视频id
			VideoId     string `json:"video_id"`     //视频真实id
			ShareUrl    string `json:"share_url"`    //视频播放页面。视频播放页可能会失效
			Title       string `json:"title"`        //视频标题
			VideoStatus int64  `json:"video_status"` //表示视频状态。1:细化为5、6、7三种状态;2:不适宜公开;4:审核中;5:公开视频;6:好友可见;7:私密视频
			Statistics  struct {
				CommentCount  int64 `json:"comment_count"`  //评论数
				DiggCount     int64 `json:"digg_count"`     //点赞数
				DownloadCount int64 `json:"download_count"` //下载数
				ForwardCount  int64 `json:"forward_count"`  //转发数
				PlayCount     int64 `json:"play_count"`     //播放数，只有作者本人可见。公开视频设为私密后，播放数也会返回0。
				ShareCount    int64 `json:"share_count"`    //分享数
			} `json:"statistics"` //统计数据
		} `json:"list"`
	} `json:"data"`
	Extra *ExtraDto `json:"extra"`
}
type ExtraDto struct {
	Description    string `json:"description"`     //错误码描述
	ErrorCode      int    `json:"error_code"`      //错误码
	LogId          string `json:"logid"`           //日志ID
	Now            int64  `json:"now"`             //毫秒级时间戳
	SubDescription string `json:"sub_description"` //子错误码描述
	SubErrorCode   int    `json:"sub_error_code"`  //子错误码
}
