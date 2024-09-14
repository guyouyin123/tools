package qexmap

import (
	"fmt"
	"testing"
	"time"
)

func TestNewExpiringMap(t *testing.T) {
	type user struct {
		IdCard string
		Name   string
	}
	jeff := &user{
		IdCard: "123456",
		Name:   "jeff",
	}

	userIdMap := NewExpiringMap()

	userIdMap.Set(jeff.IdCard, jeff, time.Second*10)
	v, ok := userIdMap.Get(jeff.IdCard)
	if ok {
		fmt.Println(v)
	}
}
