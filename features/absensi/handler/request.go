package handler

import (
	"time"

	"github.com/royanqodri/Absensi/features/absensi"
)

type AbsensiRequest struct {
	ID              string       `json:"id,omitempty"`
	OverTimeMasuk   string       `json:"overtime_masuk" form:"overtime_masuk"`
	OverTimeKeluar  string       `json:"overtime_keluar" form:"overtime_keluar"`
	JamMasuk        string       `json:"jam_masuk" form:"jam_masuk"`
	JamKeluar       string       `json:"jam_keluar" form:"jam_keluar"`
	TanggalSekarang string       `json:"tanggal_sekarang" form:"tanggal_sekarang"`
	CreatedAt       time.Time    `json:"check_in" form:"check_in"`
	UpdateAt        time.Time    `json:"check_out" form:"check_out"`
	UserID          string       `json:"user_id,omitempty"`
	User            UserResponse `json:"user,omitempty"`
}

type UserRequest struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"nama_lengkap" form:"nama_lengkap"`
}

func RequestToEntity(user AbsensiRequest) absensi.AbsensiEntity {
	return absensi.AbsensiEntity{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdateAt,
		UserID:         user.UserID,
		OverTimeMasuk:  user.OverTimeMasuk,
		OverTimePulang: user.OverTimeKeluar,
		JamMasuk:       user.JamMasuk,
		JamKeluar:      user.JamKeluar,
		User:           UserRequestToEntity(UserRequest(user.User)),
	}
}

func UserRequestToEntity(user UserRequest) absensi.UserEntity {
	return absensi.UserEntity{
		ID:   user.ID,
		Name: user.Name,
	}
}
