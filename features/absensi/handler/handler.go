package handler

import "github.com/royanqodri/Absensi/features/absensi"

type AbsensiHandler struct {
	absensiService absensi.AbsensiServiceInterface
}

func New(service absensi.AbsensiServiceInterface) *AbsensiHandler {
	return &AbsensiHandler{
		absensiService: service, // Mengganti absensiServic
	}
}

// func (handler *AbsensiHandler) Edit(c echo.Context) error {
// 	idUser, _, _ := middlewares.ExtractToken(c)
// 	// idUser := "13947f80-78b9-446f-9fe4-cb25caa4bea4"
// 	id := c.Param("id_absensi")
// 	err := handler.absensiService.Edit(idUser, id)
// 	if err != nil {
// 		return helper.InternalError(c, err.Error(), nil)
// 	}

// 	return helper.SuccessWithOutData(c, "success update absen pulang")
// }

// func (handler *AbsensiHandler) Add(c echo.Context) error {
// 	idUser, _, _ := middlewares.ExtractToken(c)
// 	// idUser := "13947f80-78b9-446f-9fe4-cb25caa4bea4"
// 	err := handler.absensiService.Add(idUser)
// 	if err != nil {
// 		return helper.InternalError(c, err.Error(), nil)
// 	}
// 	return helper.SuccessWithOutData(c, "success create absen")
// }

// func (handler *AbsensiHandler) GetAll(c echo.Context) error {
// 	var qparams absensi.QueryParams
// 	page := c.QueryParam("page")
// 	itemsPerPage := c.QueryParam("itemsPerPage")

// 	if itemsPerPage == "" {
// 		qparams.IsClassDashboard = false
// 	} else {
// 		qparams.IsClassDashboard = true
// 		itemsConv, errItem := strconv.Atoi(itemsPerPage)
// 		if errItem != nil {
// 			return helper.FailedRequest(c, "item per page not valid", nil)
// 		}
// 		qparams.ItemsPerPage = itemsConv
// 	}
// 	if page == "" {
// 		qparams.Page = 1
// 	} else {
// 		pageConv, errPage := strconv.Atoi(page)
// 		if errPage != nil {
// 			return helper.FailedRequest(c, "page not valid", nil)
// 		}
// 		qparams.Page = pageConv
// 	}

// 	// Hanya memproses pencarian dan filter berdasarkan tanggal
// 	tanggal := c.QueryParam("created_at")
// 	qparams.SerachTanggal = tanggal

// 	idUser, _, _ := middlewares.ExtractToken(c)
// 	token, errToken := usernodejs.GetTokenHandler(c)
// 	if errToken != nil {
// 		return helper.Forbidden(c, "error get token", nil)
// 	}
// 	// idUser := "13947f80-78b9-446f-9fe4-cb25caa4bea4"
// 	bol, data, err := handler.absensiService.Get(token, idUser, qparams)
// 	if err != nil {
// 		return helper.InternalError(c, err.Error(), nil)
// 	}
// 	var response []AbsensiResponse
// 	for _, value := range data {
// 		response = append(response, EntityToResponse(value))
// 	}
// 	return helper.SuccessGetAll(c, "get all absensi successfully", response, bol)
// }

// func (handler *AbsensiHandler) GetAbsensiById(c echo.Context) error {
// 	userID, _, _ := middlewares.ExtractToken(c)

// 	// Dapatkan data absensi berdasarkan ID
// 	idParam := c.Param("id_absensi")
// 	absensiResult, err := handler.absensiService.GetById(idParam)
// 	if err != nil {
// 		log.Printf("Error get detail absensi: %s", err.Error())
// 		return helper.FailedRequest(c, err.Error(), nil)
// 	}

// 	// Dapatkan data user berdasarkan ID
// 	userResult, err := usernodejs.GetByIdUser(userID)
// 	if err != nil {
// 		log.Printf("Error get detail user: %s", err.Error())
// 		return helper.FailedRequest(c, err.Error(), nil)
// 	}

// 	// Format respons sesuai dengan yang diinginkan
// 	resultResponse := EntityToResponse(absensiResult)
// 	resultResponse.User.ID = userResult.ID
// 	resultResponse.User.Name = userResult.NamaLengkap

// 	return helper.Success(c, "success read absensi", resultResponse)
// }
