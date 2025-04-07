package client

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDB(dsn string) *gorm.DB {
	open, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return open
}
