package service

import "mime/multipart"

type VideoServiceImpl interface {
	Publish(data *multipart.FileHeader, title, token string) (string, error)
	PublishList(userID uint, token string) (string, error)
	Feed()
}
type VideoService struct {
}
