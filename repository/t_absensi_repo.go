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

type TAbsensiRepository interface {
	GetById(ctx echo.Context, tx *gorm.DB, id int64) (response.TAbsensiGetResponse, error)
	GetByParamsMain(ctx echo.Context, tx *gorm.DB, request request.TAbsensiGetRequest, limit, offset int) (respData []response.TAbsensiGetResponse, totalPage int64, totalData int64, err error)
	Save(ctx echo.Context, tx *gorm.DB, req []entity.TAbsensi) error
	Update(ctx echo.Context, tx *gorm.DB, req entity.TAbsensi) error
	Delete(ctx echo.Context, tx *gorm.DB, id int64) error
}

type TAbsensiRepositoryImpl struct {
}

func NewTAbsensiRepository() TAbsensiRepository {
	return &TAbsensiRepositoryImpl{}
}

func (repo TAbsensiRepositoryImpl) getTx(tx *gorm.DB) *gorm.DB {
	if tx == nil {
		return dbConn.DBConn
	}
	return tx
}

// GetById implements TAbsensiRepository.
func (repo TAbsensiRepositoryImpl) GetById(ctx echo.Context, tx *gorm.DB, id int64) (response.TAbsensiGetResponse, error) {
	data := response.TAbsensiGetResponse{}

	result := repo.getTx(tx).
		Select("id, id_user, overtime_masuk, overtime_pulang, jam_masuk, jam_keluar, created_at, updated_at").
		Table("t_absensi").
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return response.TAbsensiGetResponse{}, nil
		}
		return response.TAbsensiGetResponse{}, result.Error
	}

	return data, nil
}

// GetByParamsMain implements TAbsensiRepository.
func (repo TAbsensiRepositoryImpl) GetByParamsMain(ctx echo.Context, tx *gorm.DB, request request.TAbsensiGetRequest, limit, offset int) (respData []response.TAbsensiGetResponse, totalPage int64, totalData int64, err error) {
	query := repo.getTx(tx).
		Select("id, id_user, overtime_masuk, overtime_pulang, jam_masuk, jam_keluar, created_at, updated_at").
		Table("t_absensi").
		Where("deleted_at IS NULL")

	if request.IdUser != 0 {
		query = query.Where("id_user = ?", request.IdUser)
	}

	if request.JamMasuk != "" {
		query = query.Where("jam_masuk = ?", request.JamMasuk)
	}

	if request.JamKeluar != "" {
		query = query.Where("jam_keluar = ?", request.JamKeluar)
	}

	if !request.StartDate.IsZero() {
		query = query.Where("DATE(created_at) >= ?", request.StartDate.Format("2006-01-02"))
	}

	if !request.EndDate.IsZero() {
		query = query.Where("DATE(created_at) <= ?", request.EndDate.Format("2006-01-02"))
	}

	if !request.CreatedAt.IsZero() {
		query = query.Where("DATE(created_at) = ?", request.CreatedAt.Format("2006-01-02"))
	}

	if !request.UpdatedAt.IsZero() {
		query = query.Where("DATE(updated_at) = ?", request.UpdatedAt.Format("2006-01-02"))
	}

	query = query.Order("created_at DESC, jam_masuk ASC")

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

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err = query.Scan(&respData).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.TAbsensiGetResponse{}, totalPage, totalData, nil
		}
		return []response.TAbsensiGetResponse{}, totalPage, totalData, err
	}

	return respData, totalPage, totalData, nil
}

// Save implements TAbsensiRepository.
func (repo *TAbsensiRepositoryImpl) Save(ctx echo.Context, tx *gorm.DB, req []entity.TAbsensi) error {
	for _, absensi := range req {
		query := "id_user = ? AND DATE(created_at) = DATE(?) AND deleted_at IS NULL"
		args := []any{
			absensi.IdUser,
			absensi.CreatedAt,
		}

		if absensi.Id != 0 {
			query += " AND id != ?"
			args = append(args, absensi.Id)
		}

		if err := repo.getTx(tx).Save(&absensi).Error; err != nil {
			return err
		}
	}

	return nil
}

// Update implements TAbsensiRepository.
func (repo *TAbsensiRepositoryImpl) Update(ctx echo.Context, tx *gorm.DB, req entity.TAbsensi) error {
	var existing entity.TAbsensi

	// Check if record exists
	err := repo.getTx(tx).Where("id = ? AND deleted_at IS NULL", req.Id).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("absensi record not found")
		}
		return err
	}

	// Check for duplicate with other records
	var duplicate entity.TAbsensi
	checkQuery := "id_user = ? AND DATE(created_at) = DATE(?) AND deleted_at IS NULL AND id != ?"
	err = repo.getTx(tx).Where(checkQuery, req.IdUser, req.CreatedAt, req.Id).First(&duplicate).Error
	if err == nil {
		return errors.New("absensi already exists for this user on this date")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	result := repo.getTx(tx).Model(&entity.TAbsensi{}).Where("id = ?", req.Id).Updates(req)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no record updated")
	}

	return nil
}

// Delete implements TAbsensiRepository.
func (repo *TAbsensiRepositoryImpl) Delete(ctx echo.Context, tx *gorm.DB, id int64) error {
	var existing entity.TAbsensi

	// Check if record exists
	err := repo.getTx(tx).Where("id = ? AND deleted_at IS NULL", id).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("absensi record not found")
		}
		return err
	}

	result := repo.getTx(tx).Where("id = ?", id).Delete(&entity.TAbsensi{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no record deleted")
	}

	return nil
}
