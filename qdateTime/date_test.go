package qdateTime

import (
	"fmt"
	"testing"
	"time"
)

func TestGetZeroTime(t *testing.T) {
	now := time.Now().Unix()
	fmt.Println(GetZeroTime(now))
}

func TestDemo(t *testing.T) {
	now := time.Now()
	startTime := now.AddDate(0, 0, -10).Unix()
	endTime := now.Unix()
	result := GetDates(startTime, endTime, 5)
	fmt.Println(result)
}
