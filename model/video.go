package model

import (
	"time"
)

type Video struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"index:create_at_idx"`
	UpdatedAt     time.Time
	DeletedAt     time.Time `gorm:"index"`
	AuthorID      uint      `db:"not null;" gorm:"index:author_idx"`
	PlayUrl       string    `db:"not null;"`
	CoverUrl      string    `db:"not null;"`
	Title         string    `db:"not null;"`
	FavoriteCount uint      `db:"default:0"`
	CommentCount  uint      `db:"default:0"`
}
