package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/royanqodri/Absensi/model/request"
	"github.com/royanqodri/Absensi/service"
	"github.com/royanqodri/Absensi/util"
	"github.com/royanqodri/Absensi/util/logging"
	"github.com/sirupsen/logrus"
)

type TAbsensiController interface {
	GetByParams(ctx echo.Context) error
	Post(ctx echo.Context) error
}

type TAbsensiControllerImpl struct {
	tAbsensiService service.TAbsensiService
}

func NewTAbsensiController(tAbsensiService service.TAbsensiService) TAbsensiController {
	return &TAbsensiControllerImpl{
		tAbsensiService: tAbsensiService,
	}
}

// GetByParams godoc
// @Summary Get Absensi
// @Description Get list of Absensi by filter and pagination
// @Tags Absensi
// @Param id_user query string false "ID User"
// @Param jam_masuk query string false "Jam Masuk"
// @Param jam_keluar query string false "Jam Keluar"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Param created_at query string false "Created At (YYYY-MM-DD)"
// @Param page_now query int false "Page Number"
// @Param limit query int false "Limit per page"
// @Param source query string true "Source"
// @Accept  json
// @Produce  json
// @Success 200 {object} response.MainResponse{data=[]response.TAbsensiGetResponse}
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /absensi [get]
func (controller *TAbsensiControllerImpl) GetByParams(ctx echo.Context) error {
	req := request.TAbsensiGetRequest{}

	// Bind query parameters
	if err := ctx.Bind(&req); err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{
			"endpoint": ctx.Request().URL.String(),
			"method":   ctx.Request().Method,
			"message":  err,
			"source":   ctx.QueryParam("source"),
		})
		util.WriteResponseEcho(ctx, util.JSON, http.StatusBadRequest, []any{}, fmt.Sprintf("%v", err), 0, 0, []string{})
		return err
	}

	// Get data from service
	respData, totalPage, totalData, err := controller.tAbsensiService.GetByParams(ctx, req)
	if err != nil {
		logging.LogWithFields(logging.ERROR, logging.ERROR, logrus.Fields{
			"endpoint": ctx.Request().URL.String(),
			"method":   ctx.Request().Method,
			"message":  err,
			"source":   ctx.QueryParam("source"),
		})
		util.WriteResponseEcho(ctx, util.JSON, http.StatusInternalServerError, []any{}, fmt.Sprintf("%v", err), 0, 0, []string{err.Error()})
		return err
	}

	util.WriteResponseEcho(ctx, util.JSON, http.StatusOK, respData, "success", totalPage, totalData, []string{})
	return err
}

// Post godoc
// @Summary Create/Update Absensi
// @Description Create or Update one or more Absensi records (bulk support)
// @Tags Absensi
// @Param source query string true "Source"
// @Accept  json
// @Produce  json
// @Param request body request.TAbsensiPostRequest true "Absensi Post Payload"
// @Success 200 {object} response.MainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 409 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /absensi [post]
func (controller *TAbsensiControllerImpl) Post(ctx echo.Context) error {
	req := request.TAbsensiPostRequest{}

	if err := ctx.Bind(&req); err != nil {

		util.WriteResponseEcho(ctx, util.JSON, http.StatusBadRequest, []any{}, fmt.Sprintf("%v", err), 0, 0, []string{})
		return err
	}

	// // Optional: Check permission (sesuaikan dengan middleware Anda)
	// if !ctx.Get("access_create").(bool) && !ctx.Get("access_update").(bool) {
	// 	util.WriteResponseEcho(ctx, util.JSON, http.StatusForbidden, nil, "Access denied. Please check your permissions.", 0, 0, []string{})
	// 	return fmt.Errorf("access denied")
	// }

	err := controller.tAbsensiService.Post(ctx, req)
	if err != nil {

		// Handle duplicate error jika ada
		var dupErr *util.DuplicateError
		if errors.As(err, &dupErr) {
			util.WriteDuplicateResponseEcho(ctx, dupErr.Fields)
			return err
		}

		util.WriteResponseEcho(ctx, util.JSON, http.StatusInternalServerError, []any{}, fmt.Sprintf("%v", err), 0, 0, []string{err.Error()})
		return err
	}

	util.WriteResponseEcho(ctx, util.JSON, http.StatusOK, []any{}, "success", 0, 0, []string{})
	return nil
}
