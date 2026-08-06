package repository

import (
	"errors"

	"github.com/labstack/echo/v4"
	dbConn "github.com/royanqodri/Absensi/database"
	"github.com/royanqodri/Absensi/model/entity"
	"github.com/royanqodri/Absensi/model/request"
	"github.com/royanqodri/Absensi/model/response"
	"gorm.io/gorm"
)

type TUserRepository interface {
	GetById(ctx echo.Context, tx *gorm.DB, id int64) (response.TUserGetResponse, error)
	GetByParamsMain(ctx echo.Context, tx *gorm.DB, request request.TUserGetRequest, limit, offset int) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error)
	Save(ctx echo.Context, tx *gorm.DB, req []entity.TUser) error
	Delete(ctx echo.Context, tx *gorm.DB, id int64) error
}

type TUserRepositoryImpl struct {
}

func NewTUserRepository() TUserRepository {
	return &TUserRepositoryImpl{}
}

func (repo TUserRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

// GetById implements TUserRepository.
func (repo TUserRepositoryImpl) GetById(ctx echo.Context, tx *gorm.DB, id int64) (response.TUserGetResponse, error) {
	data := response.TUserGetResponse{}

	result := repo.getTx(tx).
		Select("id, username, password, name, phone_number, email, address, profile_photo, upload_ktp, created_at, updated_at").
		Table("t_user").
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return response.TUserGetResponse{}, nil
		}
		return response.TUserGetResponse{}, result.Error
	}

	return data, nil
}

// GetByParamsMain implements TUserRepository.
func (repo TUserRepositoryImpl) GetByParamsMain(ctx echo.Context, tx *gorm.DB, request request.TUserGetRequest, limit, offset int) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error) {
	query := repo.getTx(tx).
		Select("id, username, password, name, phone_number, email, address, profile_photo, upload_ktp, created_at, updated_at").
		Table("t_user").
		Where("deleted_at IS NULL")

	if request.Username != "" {
		query = query.Where("username = ?", request.Username)
	}

	if request.Name != "" {
		query = query.Where("name LIKE ?", "%"+request.Name+"%")
	}

	if request.PhoneNumber != "" {
		query = query.Where("phone_number = ?", request.PhoneNumber)
	}

	if request.Email != "" {
		query = query.Where("email = ?", request.Email)
	}

	if request.Address != "" {
		query = query.Where("address LIKE ?", "%"+request.Address+"%")
	}

	if !request.LastTime.IsZero() {
		query = query.Where("DATE(updated_at) = ?", request.LastTime.Format("2006-01-02"))
	}

	// Order by
	query = query.Order("created_at DESC, name ASC")

	// Count total data
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, 0, err
	}

	// Calculate total pages
	if limit > 0 {
		totalPage = (totalData + int64(limit) - 1) / int64(limit)
	} else {
		totalPage = 1
	}

	// Apply pagination
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Execute query
	if err = query.Scan(&respData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.TUserGetResponse{}, totalPage, totalData, nil
		}
		return []response.TUserGetResponse{}, totalPage, totalData, err
	}

	return respData, totalPage, totalData, nil
}

// Save implements TUserRepository.
func (repo *TUserRepositoryImpl) Save(ctx echo.Context, tx *gorm.DB, req []entity.TUser) error {
	for _, user := range req {
		query := "username = ? AND DATE(created_at) = DATE(?) AND deleted_at IS NULL"
		args := []any{
			user.Username,
			user.CreatedAt,
		}

		if user.Id != 0 {
			query += " AND id != ?"
			args = append(args, user.Id)
		}

		if err := repo.getTx(tx).Omit("created_at").Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete implements TUserRepository.
func (repo *TUserRepositoryImpl) Delete(ctx echo.Context, tx *gorm.DB, id int64) error {
	var existing entity.TUser

	// Check if record exists
	err := repo.getTx(tx).Where("id = ? AND deleted_at IS NULL", id).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user record not found")
		}
		return err
	}

	// Soft delete
	result := repo.getTx(tx).Where("id = ?", id).Delete(&entity.TUser{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no record deleted")
	}

	return nil
}
