package util

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func readJSONFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename) // Gunakan os.ReadFile sebagai pengganti ioutil.ReadFile
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func JsonResponseHandler(c echo.Context, filename string) error {
	// Cek dulu apakah file ada
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "File tidak ditemukan",
		})
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		c.Logger().Errorf("gagal baca file %s: %v", filename, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal membaca file",
		})
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Format JSON tidak valid",
		})
	}

	return c.JSON(http.StatusOK, jsonData)
}
