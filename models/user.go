package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	TgID             int64
	LastSnapshotTime time.Time
	SaveMode         string `gorm:"default:'pdf'"`
	Busy             bool
}
