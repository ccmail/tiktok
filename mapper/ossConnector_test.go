package mapper

import (
	"fmt"
	"testing"
)

func TestInitOSS(t *testing.T) {
	err := InitOSS()
	if err != nil {
		t.Error("连接OSS失败")
	}
	fmt.Printf("%v\n", ossConfig)
	fmt.Printf("%v\n", Bucket)
	fmt.Println(BucketName)
	fmt.Println(EndPoint)
}
