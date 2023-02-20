package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	AuthorID      uint   `db:"not null;"`
	PlayUrl       string `db:"not null;"`
	CoverUrl      string `db:"not null;"`
	Title         string `db:"not null;"`
	FavoriteCount uint   `db:"default:0"`
	CommentCount  uint   `db:"default:0"`
}
