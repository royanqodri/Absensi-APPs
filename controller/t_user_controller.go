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

type TUserController interface {
	GetByParams(ctx echo.Context) error
	Post(ctx echo.Context) error
}

type TUserControllerImpl struct {
	tUserService service.TUserService
}

func NewTUserController(tUserService service.TUserService) TUserController {
	return &TUserControllerImpl{
		tUserService: tUserService,
	}
}

// GetByParams godoc
// @Summary Get User
// @Description Get list of User by filter and pagination
// @Tags User
// @Param created_at query string false "Created At (YYYY-MM-DD)"
// @Param page_now query int false "Page Number"
// @Param limit query int false "Limit per page"
// @Accept  json
// @Produce  json
// @Success 200 {object} response.MainResponse{data=[]response.TUserGetResponse}
// @Failure 400 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /user [get]
func (controller *TUserControllerImpl) GetByParams(ctx echo.Context) error {
	req := request.TUserGetRequest{}

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
	respData, totalPage, totalData, err := controller.tUserService.GetByParams(ctx, req)
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
// @Summary Create/Update User
// @Description Create or Update one or more User records (bulk support)
// @Tags User
// @Accept  json
// @Produce  json
// @Param request body request.TUserPostRequest true "User Post Payload"
// @Success 200 {object} response.MainResponse
// @Failure 400 {object} response.MainResponse
// @Failure 409 {object} response.MainResponse
// @Failure 500 {object} response.MainResponse
// @Router /user [post]
func (controller *TUserControllerImpl) Post(ctx echo.Context) error {
	req := request.TUserPostRequest{}

	if err := ctx.Bind(&req); err != nil {

		util.WriteResponseEcho(ctx, util.JSON, http.StatusBadRequest, []any{}, fmt.Sprintf("%v", err), 0, 0, []string{})
		return err
	}

	err := controller.tUserService.Post(ctx, req)
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
