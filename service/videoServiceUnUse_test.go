package service

import (
	"testing"
	"tiktok/pkg/middleware"
)

var path = "E:\\cc\\videos\\juejin\\publish测 试 文 件.mp4"

func TestExampleReadFrameAsJpeg(t *testing.T) {
	_, err := middleware.ExampleReadFrameAsJpeg(path, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestOssUploadFromPath(t *testing.T) {
	_ = middleware.InitOSSSupportTest()
	url, err := middleware.OssUploadFromPath("publi sh测 \\试文件.mp4", path)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(url)
	}
}
