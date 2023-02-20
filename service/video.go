package service

import (
	"mime/multipart"
	"tiktok/pkg/common"
)

type VideoServiceImpl interface {
	Publish(*multipart.FileHeader, string, string) (string, error)
	PublishList(uint, string) ([]common.VideoResp, error)
	Feed(string, string) ([]common.VideoResp, int64, error)
}
type VideoService struct {
}
