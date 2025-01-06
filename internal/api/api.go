package api

import (
	"encoding/json"
	"github.com/schmalz302/Calc/internal/math"
	"net/http"
	"strings"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := math.Calc(req.Expression)
	
	

	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "mismatched") {
			writeErrorResponse(w, "Invalid expression", http.StatusUnprocessableEntity)
		} else {
			writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
