package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           map[string]string
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:           "корректное выражение",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "2+2"},
			expectedStatus: http.StatusOK,
			expectedResult: 4,
			expectError:    false,
		},
		{
			name:           "корректное выражение со скообками",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "(2+2)*3"},
			expectedStatus: http.StatusOK,
			expectedResult: 12,
			expectError:    false,
		},
		{
			name:           "пустое выражение",
			method:         http.MethodPost,
			body:           map[string]string{"expression": ""},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name:           "некорректный метод",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:           "некорректное выражение",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "2++2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name:           "некорректное выражени, отсутствует закрывающая скобка",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "(2+2"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectError:    true,
		},
		{
			name:           "деление на ноль",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "1/0"},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest(tc.method, "/api/v1/calculate", bytes.NewBuffer(bodyBytes))
			} else {
				req = httptest.NewRequest(tc.method, "/api/v1/calculate", nil)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CalculateHandler)

			handler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
			}

			if !tc.expectError && tc.expectedStatus == http.StatusOK {
				var response struct {
					Result float64 `json:"result"`
				}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Could not decode response: %v", err)
				}
				if response.Result != tc.expectedResult {
					t.Errorf("handler returned wrong result: got %v want %v", response.Result, tc.expectedResult)
				}
			}
		})
	}
}