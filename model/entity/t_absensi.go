package entity

import (
	"time"

	"gorm.io/gorm"
)

type TAbsensi struct {
	Id             int64          `gorm:"primaryKey;column:id;autoIncrement;type:bigint" json:"id"`
	IdUser         int64          `gorm:"column:id_user;type:bigint;not null"`
	OverTimeMasuk  string         `gorm:"column:overtime_masuk;type:varchar(50)"`
	OverTimePulang string         `gorm:"column:overtime_pulang;type:varchar(50)"`
	JamMasuk       string         `gorm:"column:jam_masuk;type:varchar(50);not null"`
	JamKeluar      string         `gorm:"column:jam_keluar;type:varchar(50)"`
	StatusData     string         `gorm:"column:status_data;type:varchar(10);default:''"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TAbsensi) TableName() string {
	return "t_absensi"
}
