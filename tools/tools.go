package tools

import (
	"fmt"
	"strings"
	"time"
)

func ParseTime(strTime string) (int64, error) {
	//str1 := "2022-1-1"           //1.yyyy-m-d
	//str2 := "2022-01-01"         //2.yyyy-mm-dd
	//str3 := "2022/01/01"         //3:yyyy/mm/dd
	//str4 := "2022/1/1"           //4:yyyy/m/d
	//str5 := "2022-1-1 00-00-00"  //5:yyyy-mm-dd hh-mm-ss
	//str6 := "2022-1-1 00:00:00"  //6:yyyy-mm-dd hh:mm:ss
	//str7 := "2022-1-1 00/00/00"  //7:yyyy-mm-dd hh/mm/ss
	//str8 := "2022/1/1 00-00-00"  //8:yyyy/mm/dd hh-mm-ss
	//str9 := "2022/1/1 00:00:00"  //9:yyyy/mm/dd hh:mm:ss
	//str10 := "2022/1/1 00/00/00" //10:yyyy/mm/dd hh/mm/ss

	strTime = strings.TrimSpace(strTime)
	s := strings.Split(strTime, " ")
	layout := "2006-1-2 15-04-05" //默认格式
	switch len(s) {
	case 1: //年月日
		if strings.Contains(strTime, "-") {
			layout = "2006-1-2" //1,2:yyy-m-d,yyy-mm-dd
		} else {

			layout = "2006/1/2" //3,4:yyyy/mm/dd,yyyy/m/d
		}
	case 2: //年月日时分秒
		if strings.Contains(s[0], "-") {
			if strings.Contains(s[1], "-") {
				layout = "2006-1-2 15-04-05" //5:yyyy-mm-dd hh-mm-ss
			} else {
				if strings.Contains(s[1], ":") {
					layout = "2006-1-2 15:04:05" //6:yyyy-mm-dd hh:mm:ss
				} else {
					layout = "2006-1-2 15/04/05" //7:yyyy-mm-dd hh/mm/ss
				}
			}
		} else {
			if strings.Contains(s[1], "-") {
				layout = "2006/1/2 15-04-05" //8:yyyy/mm/dd hh-mm-ss
			} else {
				if strings.Contains(s[1], ":") {
					layout = "2006/1/2 15:04:05" //9:yyyy/mm/dd hh:mm:ss
				} else {
					layout = "2006/1/2 15/04/05" //10:yyyy/mm/dd hh/mm/ss
				}
			}
		}
	}
	shipTsInt, err := time.ParseInLocation(layout, strTime, time.Local)
	if err != nil {
		fmt.Println("err:", err)
		return 0, err
	}
	return shipTsInt.Unix(), nil
}
