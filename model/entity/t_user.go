package entity

import (
	"time"

	"gorm.io/gorm"
)

type TUser struct {
	Id           int64          `gorm:"primaryKey;column:id;autoIncrement;type:bigint" json:"id"`
	Name         string         `gorm:"column:name;type:varchar(100);not null"`
	Username     string         `gorm:"column:username;type:varchar(250);not null"`
	Email        string         `gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
	PhoneNumber  string         `gorm:"column:phone_number;type:varchar(20);uniqueIndex"`
	Password     string         `gorm:"column:password;type:varchar(255);not null"`
	Address      string         `gorm:"column:address;type:text"`
	ProfilePhoto string         `gorm:"column:profile_photo;type:varchar(500)"`
	UploadKtp    string         `gorm:"column:upload_ktp;type:varchar(255)"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TUser) TableName() string {
	return "t_user"
}
