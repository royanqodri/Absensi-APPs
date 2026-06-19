package entity

import (
	"time"

	"gorm.io/gorm"
)

type TAbsensi struct {
	Id             string `gorm:"primaryKey"`
	IdUser         string
	OverTimeMasuk  string
	OverTimePulang string
	JamMasuk       string
	JamKeluar      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (b *TAbsensi) TableName() string {
	return "t_absensi"
}
