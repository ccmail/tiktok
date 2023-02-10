package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	AuthorID      uint   `gorm:"not null;"`
	PlayUrl       string `gorm:"not null;"`
	CoverUrl      string `gorm:"not null;"`
	Title         string `gorm:"not null;"`
	FavoriteCount uint
	CommentCount  uint
}
