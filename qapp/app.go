package qapp

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type Results struct {
	TrackName     string   `json:"trackName"`     //名字
	Genres        []string `json:"genres"`        //标签
	ArtworkUrl512 string   `json:"artworkUrl512"` //图片
	BundleId      string   `json:"bundleId"`      //包名
}

type ItunesApple struct {
	ResultCount uint32     `json:"resultCount"`
	Results     []*Results `json:"results"`
}

type GoogleData struct {
	Name string `json:"name"`
	Img  string `json:"img"`
}

// 获取app应用信息
func GetItunesAppleUrl(id, cn string) (*ItunesApple, error) {
	/*
		cn: http://itunes.apple.com/cn/lookup?id=appId (appId为应用 id)
		en: http://itunes.apple.com/lookup?id=appId (appId为应用 id)
	*/
	boo := 0
start:
	boo++
	url := fmt.Sprintf("http://itunes.apple.com/lookup?id=%s", id)
	if cn == "cn" {
		url = fmt.Sprintf("http://itunes.apple.com/cn/lookup?id=%s", id)
	}
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	info := &ItunesApple{}
	_ = jsoniter.Unmarshal(body, &info)
	if len(info.Results) == 0 {
		if boo <= 1 {
			cn = "cn"
			goto start
		}
		return nil, nil
	}
	return info, nil
}

// 获取Google应用信息
func GetGoogleUrl(id string) (*GoogleData, error) {
	/*
		url = "https://play.google.com/store/apps/details?id=com.google.android.youtube"
	*/
	url := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s", id)
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	buf := string(body)
	info := &GoogleData{}

	reName := regexp.MustCompile(`itemprop="name">(.*?)</h1>`)
	reNameList := reName.FindStringSubmatch(buf)
	if len(reNameList) > 1 {
		name := reNameList[1]
		info.Name = name
	}
	reImg := regexp.MustCompile(`og:image" content="(.*?)">`)
	reImgList := reImg.FindStringSubmatch(buf)
	if len(reImgList) > 1 {
		img := reImgList[1]
		info.Img = img
	}
	return info, nil
}

// 检查app分享链接类型
func CheckUrl(url string) (string, string, error) {
	/*
		支持3中类型url:apps.apple.com，apps.apple.com/cn，play.google.com
			url := "https://apps.apple.com/app/%E6%8A%96%E9%9F%B3/id1142110895"
			url := "https://apps.apple.com/cn/app/%E6%8A%96%E9%9F%B3/id1142110895"
			url := "https://apps.apple.com/cn/app/tencent-meeting/id1484048379?l=en-GB"
			uel := "https://apps.apple.com/cn/app/%E8%85%BE%E8%AE%AF%E4%BC%9A%E8%AE%AE-%E5%A4%9A%E4%BA%BA%E5%AE%9E%E6%97%B6%E8%A7%86%E9%A2%91%E4%BC%9A%E8%AE%AE%E8%BD%AF%E4%BB%B6/id1484048379"

			url := "https://play.google.com/store/apps/details?id=com.google.android.youtube"
	*/

	appleCn := strings.Contains(url, "apple.com/cn")
	id := ""
	zh := ""
	if appleCn {
		sp := strings.Split(url, "id")
		if len(sp) < 1 {
			return "", "", fmt.Errorf("url fail")
		}
		id = sp[1]
		zh = "cn"
		sp2 := strings.Split(id, "?")
		if len(sp2) > 0 {
			id = sp2[0]
		}
		return id, zh, nil
	}

	appleEn := strings.Contains(url, "apple.com")
	if appleEn {
		sp := strings.Split(url, "id")
		if len(sp) < 1 {
			return "", "", fmt.Errorf("url fail")
		}
		id = sp[1]
		zh = "en"
		sp2 := strings.Split(id, "?")
		if len(sp2) > 0 {
			id = sp2[0]
		}
		return id, zh, nil
	}
	google := strings.Contains(url, "play.google.com")
	if google {
		sp := strings.Split(url, "id=")
		if len(sp) < 1 {
			return "", "", fmt.Errorf("url fail")
		}
		id = sp[1]
		zh = "google"
		return id, zh, nil
	}
	return id, zh, nil
}
