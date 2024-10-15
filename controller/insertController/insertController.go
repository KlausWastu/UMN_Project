package insertcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/jeypc/go-jwt-mux/helper"
	"github.com/jeypc/go-jwt-mux/models"
)

func InsertMultiple(w http.ResponseWriter, r *http.Request) {
	// menagambil inputan json
	var report models.RequestInsertReport
	var tempData []models.TemporaryData
	var reportScores []models.ReportScore
	// mendapat json dari body
	decoderJSON := json.NewDecoder(r.Body)
	err := decoderJSON.Decode(&report)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// Cek apakah ada data yang akan di-insert
	if len(report.Value) == 0 {
		response := map[string]string{"message": "No data provided"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	err = models.DB.Where("no IN ?", report.Value).Find(&tempData).Error
	if err != nil {
		response := map[string]string{"message": "Email or password not found"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// Mapping data dari TemporaryData ke ReportScore
	for _, temp := range tempData {
		reportScore := models.ReportScore{
			Name:  temp.Name,
			Score: temp.Score,
		}
		reportScores = append(reportScores, reportScore)
	}

	// Insert multiple records ke tabel ReportScore menggunakan GORM
	err = models.DB.Create(&reportScores).Error
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// Hapus semua data dari TemporaryData setelah berhasil insert ke ReportScore
	err = models.DB.Where("true").Delete(&models.TemporaryData{}).Error
	if err != nil {
		response := map[string]string{"message": "Failed to delete temporary data"}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}
