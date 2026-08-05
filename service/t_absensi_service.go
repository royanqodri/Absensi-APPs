package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	dbConn "github.com/royanqodri/Absensi/database"
	"github.com/royanqodri/Absensi/model/entity"
	"github.com/royanqodri/Absensi/model/request"
	"github.com/royanqodri/Absensi/model/response"
	"github.com/royanqodri/Absensi/repository"
	"github.com/royanqodri/Absensi/util"
	"github.com/royanqodri/Absensi/util/constants"
)

type TAbsensiService interface {
	GetByParams(ctx echo.Context, request request.TAbsensiGetRequest) (respData []response.TAbsensiGetResponse, totalPage int64, totalData int64, err error)
	Post(ctx echo.Context, request request.TAbsensiPostRequest) (err error) // insert, update, delete
}

type TAbsensiServiceImpl struct {
	tAbsensiRepo repository.TAbsensiRepository
}

func NewTAbsensiService(tAbsensiRepo repository.TAbsensiRepository) TAbsensiService {
	return &TAbsensiServiceImpl{
		tAbsensiRepo: tAbsensiRepo,
	}
}

// GetByParams implements TAbsensiService.
func (service TAbsensiServiceImpl) GetByParams(ctx echo.Context, request request.TAbsensiGetRequest) (respData []response.TAbsensiGetResponse, totalPage int64, totalData int64, err error) {
	// Parse pagination parameters (Echo style)
	pageStr := ctx.QueryParam("page_now")
	limitStr := ctx.QueryParam("limit")

	// Default values
	var pageNum int64 = 1
	var limitNum int64 = 0
	var offset int64 = 0

	// Parse page_now
	if pageStr != "" {
		pageNum, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || pageNum < 1 {
			pageNum = 1
		}
	}

	// Parse limit
	if limitStr != "" {
		limitNum, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limitNum < 1 {
			limitNum = 10
		}
	} else if pageStr != "" {
		// Jika hanya page_now yang dikirim, gunakan default limit 10
		limitNum = 10
	}

	if limitNum > 0 {
		offset = (pageNum - 1) * limitNum
	}

	// Panggil Repository
	respData, totalPage, totalData, err = service.tAbsensiRepo.GetByParamsMain(ctx, nil, request, int(limitNum), int(offset))
	if err != nil {
		return []response.TAbsensiGetResponse{}, 0, 0, fmt.Errorf("failed to get paginated data: %w", err)
	}

	// Hitung totalPage
	if limitNum > 0 && totalData > 0 {
		totalPage = (totalData + limitNum - 1) / limitNum
	} else {
		totalPage = 1
	}

	if totalData == 0 {
		return []response.TAbsensiGetResponse{}, totalPage, totalData, nil
	}

	return respData, totalPage, totalData, nil
}

// Post for Insert & Update
func (service TAbsensiServiceImpl) Post(ctx echo.Context, request request.TAbsensiPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entity.TAbsensi, len(request.Data))

	for i, data := range request.Data {
		// Trim data
		operation := constants.STATUS_DATA_INSERT
		if data.Id != 0 {
			operation = constants.STATUS_DATA_UPDATE
		}

		if operation == constants.STATUS_DATA_INSERT {
			dataReq[i] = entity.TAbsensi{
				IdUser:         data.IdUser,
				OverTimeMasuk:  data.OverTimeMasuk,
				OverTimePulang: data.OverTimePulang,
				JamMasuk:       data.JamMasuk,
				JamKeluar:      data.JamKeluar,
				CreatedAt:      util.GetTimeNowByLoc(),
				UpdatedAt:      util.GetTimeNowByLoc(),
			}
		} else { // UPDATE
			existing, errGet := service.tAbsensiRepo.GetById(ctx, tx, data.Id)
			if errGet != nil {

				return errGet
			}
			if existing.Id == "" {
				return errors.New("absensi record not found")
			}

			dataReq[i] = entity.TAbsensi{
				Id:             data.Id,
				IdUser:         data.IdUser,
				OverTimeMasuk:  data.OverTimeMasuk,
				OverTimePulang: data.OverTimePulang,
				JamMasuk:       data.JamMasuk,
				JamKeluar:      data.JamKeluar,
				CreatedAt:      existing.CreatedAt,
				UpdatedAt:      util.GetTimeNowByLoc(),
			}
		}

	}

	if len(dataReq) > 0 {
		if err = service.tAbsensiRepo.Save(ctx, tx, dataReq); err != nil {

			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
