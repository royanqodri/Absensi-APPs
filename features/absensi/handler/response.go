package handler

import (
	"github.com/royanqodri/Absensi/features/absensi"
)

type AbsensiResponse struct {
	ID              string       `json:"id,omitempty"`
	OverTimeMasuk   string       `json:"overtime_masuk" form:"overtime_masuk"`
	OverTimeKeluar  string       `json:"overtime_keluar" form:"overtime_keluar"`
	JamMasuk        string       `json:"jam_masuk" form:"jam_masuk"`
	JamKeluar       string       `json:"jam_keluar" form:"jam_keluar"`
	TanggalSekarang string       `json:"tanggal_sekarang" form:"tanggal_sekarang"`
	CreatedAt       string       `json:"check_in" form:"check_in"`
	UpdateAt        string       `json:"check_out" form:"check_out"`
	UserID          string       `json:"user_id,omitempty"`
	User            UserResponse `json:"user,omitempty"`
}

type UserResponse struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"nama_lengkap,omitempty"`
}

func EntityToResponse(user absensi.AbsensiEntity) AbsensiResponse {

	return AbsensiResponse{
		ID:              user.ID,
		OverTimeMasuk:   user.OverTimeMasuk,
		OverTimeKeluar:  user.OverTimePulang,
		JamMasuk:        user.JamMasuk,
		JamKeluar:       user.JamKeluar,
		TanggalSekarang: user.CreatedAt.Format("2006-01-02"),
		CreatedAt:       user.CreatedAt.Format("15:04:05.000"),
		UpdateAt:        user.UpdatedAt.Format("15:04:05.000"),
		UserID:          user.UserID,
		User:            UserEntityToResponse(user.User),
	}
}

func UserEntityToResponse(user absensi.UserEntity) UserResponse {
	return UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}
