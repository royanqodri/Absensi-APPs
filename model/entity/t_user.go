package entity

import (
	"time"

	"gorm.io/gorm"
)

type Tuser struct {
	Id           int64
	Name         string
	Email        string
	PhoneNumber  string
	Password     string
	Address      string
	ProfilePhoto string
	UploadKtp    string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (b *Tuser) TableName() string {
	return "t_user"
}
