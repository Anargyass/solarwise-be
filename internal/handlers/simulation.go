package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"solar-backend/internal/models"
	"solar-backend/internal/services"
)

type errorResponse struct {
	Error string `json:"error"`
}

func SimulateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.SimulationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Lokasi = strings.TrimSpace(req.Lokasi)
	if req.Lokasi == "" || req.TagihanBulanan <= 0 {
		writeJSONError(w, http.StatusBadRequest, "lokasi and tagihan_bulanan (> 0) are required")
		return
	}

	resp, err := services.CalculateFeasibility(req)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "required") ||
			strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "Data not available for this location") ||
			strings.Contains(errMsg, "invalid") {
			writeJSONError(w, http.StatusBadRequest, errMsg)
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}
