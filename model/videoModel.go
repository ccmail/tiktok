package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	AuthorID      uint64 `gorm:"not null;"`
	PlayURL       string `gorm:"not null;"`
	CoverUrl      string `gorm:"not null;"`
	Title         string `gorm:"not null;"`
	FavoriteCount string
	CommentCount  string
}
