package request

import "time"

type TAbsensiGetRequest struct {
	IdUser    int64     `form:"id_user"`
	StartDate time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" time_format:"2006-01-02"`
	JamMasuk  string    `form:"jam_masuk"`
	JamKeluar string    `form:"jam_keluar"`
	CreatedAt time.Time `form:"created_at" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `form:"updated_at" time_format:"2006-01-02 15:04:05"`
}

type TAbsensiPostRequest struct {
	Data []TAbsensiPostDetailRequest `json:"data" binding:"required"`
}

type TAbsensiPostDetailRequest struct {
	Type           string `json:"type"`
	Id             int64  `json:"id"`
	IdUser         int64  `json:"id_user" binding:"required"`
	OverTimeMasuk  string `json:"overtime_masuk"`
	OverTimePulang string `json:"overtime_pulang"`
	JamMasuk       string `json:"jam_masuk" binding:"required"`
	JamKeluar      string `json:"jam_keluar"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type TAbsensiDeleteRequest struct {
	Id string `json:"id" binding:"required"`
}
