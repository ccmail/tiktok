package service

import (
	"testing"
	"tiktok/mapper"
)

var path = "E:\\cc\\videos\\juejin\\publish测试文件.mp4"

func TestExampleReadFrameAsJpeg(t *testing.T) {
	_, err := ExampleReadFrameAsJpeg(path, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestOssUploadFromPath(t *testing.T) {
	mapper.InitOSS()
	url, err := OssUploadFromPath("publish测试文件.mp4", path)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(url)
	}
}
