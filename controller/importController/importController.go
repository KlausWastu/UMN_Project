package importController

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func ImportFile(w http.ResponseWriter, r *http.Request) {
	// mengambil file dari body postman
	file, _, err := r.FormFile("excel")
	if err != nil {
		response := map[string]string{"message": "failed to load excel"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer file.Close()

	// menyimpan file sementara
	tempFile, err := io.ReadAll(file)
	if err != nil {
		response := map[string]string{"message": "failed to read excel"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// membuka file excel dengan excelize
	excelFile, err := excelize.OpenReader(bytes.NewReader(tempFile))
	if err != nil {
		response := map[string]string{"message": "failed to open excel"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// mendapatkan daftar nama sheet
	sheetNames := excelFile.GetSheetList()
	if len(sheetNames) == 0 {
		response := map[string]string{"message": "no sheets found in the excel file"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	sheetName := sheetNames[0]
	fmt.Println("Available Sheets:", sheetNames)
	fmt.Println("Sheet:", sheetName)

	// mendapatkan baris excel dari sheet yang dipilih
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		response := map[string]string{"message": "failed to get rows"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// proses untuk menyimpan data ke struct TemporaryData
	// var data []models.TemporaryData
	var cleanData []models.TemporaryData

	sameNames := make(map[string]struct{})
	for i, row := range rows {
		if i == 0 { // skip header baris
			continue
		}
		if len(row) < 3 { // memastikan panjang data pada baris
			continue
		}

		// parsing nilai
		no, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			continue
		}
		name := row[1]
		score, err := strconv.ParseInt(row[2], 10, 8)
		if err != nil {
			continue
		}

		var existingReport models.ReportScore
		err = models.DB.Where("name = ?", name).First(&existingReport).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			// handle error selain record tidak ditemukan
			log.Println("Database error:", err)
			continue

		}

		if err == gorm.ErrRecordNotFound {
			if _, exists := sameNames[name]; !exists {
				sameNames[name] = struct{}{}

				// membuat objek temporarydata dan masukan ke slicing data
				tempData := models.TemporaryData{
					No:    no,
					Name:  name,
					Score: int8(score),
				}
				// data = append(data, tempData)
				cleanData = append(cleanData, tempData)
			}
		}
	}

	// Insert multiple records ke database
	if len(cleanData) > 0 {
		err = models.DB.Create(&cleanData).Error
		if err != nil {
			response := map[string]string{"message": "failed to insert data"}
			helper.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	response := map[string]interface{}{
		"message": "import data successfully",
		// "rawData":   data,
		"cleanData": cleanData,
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}
