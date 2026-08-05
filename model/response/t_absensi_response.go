package response

import "time"

type TAbsensiGetResponse struct {
	Id             string    `json:"id"`
	IdUser         string    `json:"id_user"`
	OverTimeMasuk  string    `json:"overtime_masuk"`
	OverTimePulang string    `json:"overtime_pulang"`
	JamMasuk       string    `json:"jam_masuk"`
	JamKeluar      string    `json:"jam_keluar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TAbsensiPostResponse struct {
	Id         string `json:"id"`
	IdUser     string `json:"id_user"`
	JamMasuk   string `json:"jam_masuk"`
	JamKeluar  string `json:"jam_keluar"`
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}
