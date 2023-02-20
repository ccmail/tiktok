package db

import (
	"fmt"
	"log"
	"testing"
	"tiktok/mapper"
	"time"
)

func TestGetMessageList(t *testing.T) {
	mapper.InitDBConnectorSupportTest()
	parse, err := time.ParseInLocation("2006-01-02 15:04:05", "2023-02-15 14:13:11.041", time.Local)
	if err != nil {
		log.Println("时间解析失败")
		t.Error()
	}
	log.Println("时间解析为", parse.Unix(), parse.Nanosecond())
	fmt.Println(time.Unix(parse.Unix(), 0))
	fmt.Println(parse)
	list, err := GetMessageList(1, 2, parse)
	if err != nil {
		t.Error()
	}
	for _, v := range list {
		fmt.Println(v)
	}

}

func TestTime(t *testing.T) {
	//var a int64 = 1676902325
	milli := time.Now().UnixMilli()
	unix := time.UnixMilli(milli)
	fmt.Println(milli, unix)
}
