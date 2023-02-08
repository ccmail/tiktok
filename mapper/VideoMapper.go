package mapper

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	VideoID     uint64 `gorm:"primaryKey"`
	AuthorID    uint64
	VideoURL    string
	CoverPicURL string
	Title       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
