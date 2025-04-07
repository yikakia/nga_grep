package model

import (
	"time"
)

type ThreadLatestData struct {
	TID            int       `gorm:"column:tid;primaryKey"`
	LastTime       time.Time `gorm:"column:last_time"`
	LastReplyCount int       `gorm:"column:last_reply_count"`
}

type ThreadCount struct {
	DateTime int64 `gorm:"column:date_time;primaryKey"`
	Count    int   `gorm:"column:count"`
}
