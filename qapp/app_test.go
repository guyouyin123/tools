package qapp

import (
	"fmt"
	"testing"
)

func TestCheckUrl(t *testing.T) {
	url := "https://apps.apple.com/app/id1559181149"
	id, source, err := CheckUrl(url)
	if err != nil {
		fmt.Println(err)
	}
	switch source {
	case "en", "cn":
		info, _ := GetItunesAppleUrl(id, source)
		fmt.Println(info)
	case "google":
		googleData, _ := GetGoogleUrl(id)
		fmt.Println(googleData)
	}
}
