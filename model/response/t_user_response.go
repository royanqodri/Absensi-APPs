package response

import "time"

type TUserGetResponse struct {
	Id           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	ProfilePhoto string    `json:"profile_photo"`
	UploadKtp    string    `json:"upload_ktp"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TUserPostResponse struct {
	Id           int64     `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Address      string    `json:"address"`
	ProfilePhoto string    `json:"profile_photo"`
	UploadKtp    string    `json:"upload_ktp"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	StatusCode   int       `json:"status_code"`
	StatusText   string    `json:"status_text"`
	Message      string    `json:"message"`
}

type TUserGetMainResponse struct {
	TotalPage int64              `json:"total_page"`
	TotalData int64              `json:"total_data"`
	Data      []TUserGetResponse `json:"data"`
}
