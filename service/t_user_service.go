package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	dbConn "github.com/royanqodri/Absensi/database"
	"github.com/royanqodri/Absensi/model/entity"
	"github.com/royanqodri/Absensi/model/request"
	"github.com/royanqodri/Absensi/model/response"
	"github.com/royanqodri/Absensi/repository"
	"github.com/royanqodri/Absensi/util"
	"github.com/royanqodri/Absensi/util/constants"
)

type TUserService interface {
	GetByParams(ctx echo.Context, request request.TUserGetRequest) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error)
	Post(ctx echo.Context, request request.TUserPostRequest) (err error) // insert, update, delete
}

type TUserServiceImpl struct {
	tUserRepo repository.TUserRepository
}

func NewTUserService(tUserRepo repository.TUserRepository) TUserService {
	return &TUserServiceImpl{
		tUserRepo: tUserRepo,
	}
}

// GetByParams implements TUserService.
func (service TUserServiceImpl) GetByParams(ctx echo.Context, request request.TUserGetRequest) (respData []response.TUserGetResponse, totalPage int64, totalData int64, err error) {
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
	respData, totalPage, totalData, err = service.tUserRepo.GetByParamsMain(ctx, nil, request, int(limitNum), int(offset))
	if err != nil {
		return []response.TUserGetResponse{}, 0, 0, fmt.Errorf("failed to get paginated data: %w", err)
	}

	// Hitung totalPage
	if limitNum > 0 && totalData > 0 {
		totalPage = (totalData + limitNum - 1) / limitNum
	} else {
		totalPage = 1
	}

	if totalData == 0 {
		return []response.TUserGetResponse{}, totalPage, totalData, nil
	}

	return respData, totalPage, totalData, nil
}

// Post for Insert & Update
func (service TUserServiceImpl) Post(ctx echo.Context, request request.TUserPostRequest) (err error) {
	tx := dbConn.DBConn.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataReq := make([]entity.TUser, len(request.Data))

	for i, data := range request.Data {
		operation := constants.STATUS_DATA_INSERT
		if data.Id != 0 {
			operation = constants.STATUS_DATA_UPDATE
		}

		if operation == constants.STATUS_DATA_INSERT {
			// INSERT password
			passwordHash, err := util.HashPasswordWithMD5(data.Password)
			if err != nil {
				return err
			}

			dataReq[i] = entity.TUser{
				Name:         data.Name,
				Username:     data.Username,
				Email:        data.Email,
				PhoneNumber:  data.PhoneNumber,
				Password:     passwordHash,
				Address:      data.Address,
				ProfilePhoto: data.ProfilePhoto,
				UploadKtp:    data.UploadKtp,
				CreatedAt:    util.GetTimeNowByLoc(),
				UpdatedAt:    util.GetTimeNowByLoc(),
			}
		} else { // UPDATE
			existing, errGet := service.tUserRepo.GetById(ctx, tx, data.Id)
			if errGet != nil {

				return errGet
			}

			passwordHash := existing.Password
			if data.Password != "" && strings.TrimSpace(data.Password) != "" {
				if data.OldPassword == "" || strings.TrimSpace(data.OldPassword) == "" {
					return errors.New("old_password is required when updating password")
				}
				isPasswordMatched := util.ComparePasswordWithMD5(data.OldPassword, existing.Password)
				if !isPasswordMatched {
					return errors.New("old_password is incorrect")
				}

				passwordUpdate, err := util.HashPasswordWithMD5(data.Password)
				if err != nil {
					return err
				}
				passwordHash = passwordUpdate
			}

			if existing.Id == 0 {
				return errors.New("absensi record not found")
			}

			dataReq[i] = entity.TUser{
				Id:           data.Id,
				Name:         data.Name,
				Username:     data.Username,
				Email:        data.Email,
				PhoneNumber:  data.PhoneNumber,
				Password:     passwordHash,
				Address:      data.Address,
				ProfilePhoto: data.ProfilePhoto,
				UploadKtp:    data.UploadKtp,
				CreatedAt:    existing.CreatedAt,
				UpdatedAt:    util.GetTimeNowByLoc(),
			}
		}

	}

	if len(dataReq) > 0 {
		if err = service.tUserRepo.Save(ctx, tx, dataReq); err != nil {

			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
