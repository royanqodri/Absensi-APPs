package service

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/royanqodri/Absensi/features/absensi"
)

type AbsensiService struct {
	absensiService absensi.AbsensiDataInterface
	validate       *validator.Validate
}

// GetById implements absensi.AbsensiServiceInterface
func (service *AbsensiService) GetById(absensiID string) (absensi.AbsensiEntity, error) {
	result, err := service.absensiService.SelectById(absensiID)
	if err != nil {
		return absensi.AbsensiEntity{}, err
	}
	user, errUser := service.absensiService.SelectUserById(result.UserID)
	if errUser != nil {
		return absensi.AbsensiEntity{}, err
	}
	result.User.ID = user.ID
	result.User.Name = user.NamaLengkap
	return result, nil
}

// Add implements absensi.AbsensiServiceInterface
func (service *AbsensiService) Add(idUser string) error {
	var input absensi.AbsensiEntity

	sekarang := time.Now()

	tanggalWaktuBaru := time.Date(
		sekarang.Year(),
		sekarang.Month(),
		sekarang.Day(),
		sekarang.Hour(),
		sekarang.Minute(),
		sekarang.Second(),
		sekarang.Nanosecond(),
		time.UTC,
	)
	jamSeharusnya := "08:00:00"
	format := "15:04:05"

	t, err := time.Parse(format, jamSeharusnya)
	if err != nil {
		return errors.New("gagal parsing waktu")
	}
	waktuSeharusnyaMasuk := time.Date(sekarang.Year(), sekarang.Month(), sekarang.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	selisih := tanggalWaktuBaru.Sub(waktuSeharusnyaMasuk)
	keterlambatan := selisih.Minutes()

	keterlambatanInt := int(keterlambatan)
	konvKeterlambatan := strconv.Itoa(keterlambatanInt)

	input.JamMasuk = jamSeharusnya
	input.OverTimeMasuk = konvKeterlambatan
	input.UserID = idUser
	errInsert := service.absensiService.Insert(input)
	if errInsert != nil {
		return errInsert
	}
	return nil
}

// Edit implements absensi.AbsensiServiceInterface
func (service *AbsensiService) Edit(idUser string, id string) error {
	var input absensi.AbsensiEntity
	sekarang := time.Now()

	tanggalWaktuBaru := time.Date(
		sekarang.Year(),
		sekarang.Month(),
		sekarang.Day(),
		sekarang.Hour(),
		sekarang.Minute(),
		sekarang.Second(),
		sekarang.Nanosecond(),
		time.UTC,
	)
	jamSeharusnya := "17:00:00"
	format := "15:04:05"

	t, err := time.Parse(format, jamSeharusnya)
	if err != nil {
		return errors.New("gagal parsing waktu")
	}
	waktuSeharusnyaMasuk := time.Date(sekarang.Year(), sekarang.Month(), sekarang.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	selisih := tanggalWaktuBaru.Sub(waktuSeharusnyaMasuk)
	keterlambatan := selisih.Minutes()

	keterlambatanInt := int(keterlambatan)
	konvKeterlambatan := strconv.Itoa(keterlambatanInt)

	input.JamKeluar = jamSeharusnya
	input.OverTimePulang = konvKeterlambatan
	errUpdate := service.absensiService.Update(input, idUser, id)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}

// Get implements absensi.AbsensiServiceInterface
func (service *AbsensiService) Get(token string, idUser string, param absensi.QueryParams) (bool, []absensi.AbsensiEntity, error) {
	var totalPage int64
	nextPage := true

	// Get user's role
	user, err := service.absensiService.SelectUserById(idUser)
	if err != nil {
		log.Printf("Error getting user details: %s", err.Error())
		return false, nil, err
	}
	var data []absensi.AbsensiEntity
	if user.Jabatan == "karyawan" {
		// Karyawan can only view their own absensis
		count, karyawanData, err := service.absensiService.SelectAllKaryawan(idUser, param)
		if err != nil {
			log.Printf("Error selecting all absensis: %s", err.Error())
			return false, nil, err
		}

		if count == 0 {
			nextPage = false
		}
		data = karyawanData
		if param.IsClassDashboard {
			totalPage = count / int64(param.ItemsPerPage)
			if count%int64(param.ItemsPerPage) != 0 {
				totalPage += 1
			}

			if param.Page == int(totalPage) {
				nextPage = false
			}
			if data == nil {
				nextPage = false
			}
		}
	} else {
		count, allData, err := service.absensiService.SelectAll(token, param)
		if err != nil {
			log.Printf("Error selecting all absensis: %s", err.Error())
			return false, nil, err
		}
		if count == 0 {
			nextPage = false
		}
		data = allData
		if param.IsClassDashboard {
			totalPage = count / int64(param.ItemsPerPage)
			if count%int64(param.ItemsPerPage) != 0 {
				totalPage += 1
			}

			if param.Page == int(totalPage) {
				nextPage = false
			}
			if data == nil {
				nextPage = false
			}
		}
		log.Println("Absensis read successfully")
		return nextPage, data, nil
	}
	return nextPage, data, nil
}

// func New(service absensi.AbsensiDataInterface) absensi.AbsensiServiceInterface {
// 	return &AbsensiService{
// 		absensiService: service,
// 		validate:       validator.New(),
// 	}
// }
