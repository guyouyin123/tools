package qhttp

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	netUrl "net/url"
	"strings"
)

func Get(url string, params, header map[string]interface{}) (resp []byte, err error) {
	urlParse, _ := netUrl.Parse(url)
	p := urlParse.Query()
	for k, v := range params {
		p.Set(k, fmt.Sprintf("%v", v))
	}
	urlParse.RawQuery = p.Encode()
	url = urlParse.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	for k, v := range header {
		req.Header.Add(k, v.(string))
	}
	count := 0
start:
	count++
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && count > 5 {
		goto start
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Post(url string, params, header, data map[string]interface{}) (resp []byte, err error) {
	urlParse, _ := netUrl.Parse(url)
	p := urlParse.Query()
	for k, v := range params {
		p.Set(k, fmt.Sprintf("%v", v))
	}
	urlParse.RawQuery = p.Encode()
	url = urlParse.String()

	client := &http.Client{}
	dataS, _ := jsoniter.Marshal(data)
	dataR := strings.NewReader(string(dataS))
	req, err := http.NewRequest("POST", url, dataR)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v.(string))
	}
	count := 0
start:
	count++

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && count > 5 {
		goto start
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
