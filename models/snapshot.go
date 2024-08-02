package models

import (
	"gorm.io/gorm"
)

type Snapshot struct {
	gorm.Model
	Initiator    uint // User.ID
	Site         string
	ResultFileID string
}
