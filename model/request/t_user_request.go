package request

import "time"

type TUserGetRequest struct {
	Username    string    `form:"username"`
	Name        string    `form:"name"`
	PhoneNumber string    `form:"phone_number"`
	Email       string    `form:"email"`
	Address     string    `form:"address"`
	LastTime    time.Time `form:"last_time" time_format:"2006-01-02 15:04:05"`
}

type TUserPostRequest struct {
	Data []TUserPostDetailRequest `json:"data" binding:"required"`
}

type TUserPostDetailRequest struct {
	Id           int64  `json:"id"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	OldPassword  string `json:"old_password" binding:"required"`
	Name         string `json:"name" binding:"required"`
	PhoneNumber  string `json:"phone_number" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Address      string `json:"address" binding:"required"`
	ProfilePhoto string `json:"profile_photo"`
	UploadKtp    string `json:"upload_ktp"`
}
