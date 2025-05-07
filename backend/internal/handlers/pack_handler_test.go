package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"pack-sizer/internal/services"
)

func setupTestRouter() (*gin.Engine, *PackHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	packSizesManager := services.NewPackSizesManager()
	handler := NewPackHandler(packSizesManager)
	return router, handler
}

func TestGetPackSizes(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/pack-sizes", handler.GetPackSizes)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/pack-sizes", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "sizes")
	assert.Equal(t, []int{250, 500, 1000, 2000, 5000}, response["sizes"])
}

func TestAddPackSize(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/pack-sizes/add", handler.AddPackSize)

	tests := []struct {
		name           string
		payload        map[string]int
		expectedStatus int
	}{
		{
			name:           "Valid size",
			payload:        map[string]int{"size": 300},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Zero size",
			payload:        map[string]int{"size": 0},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative size",
			payload:        map[string]int{"size": -100},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/pack-sizes/add", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRemovePackSize(t *testing.T) {
	router, handler := setupTestRouter()
	router.DELETE("/pack-sizes/:size", handler.RemovePackSize)

	tests := []struct {
		name           string
		size           string
		expectedStatus int
	}{
		{
			name:           "Valid size",
			size:           "250",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid size format",
			size:           "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/pack-sizes/"+tt.size, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdatePackSize(t *testing.T) {
	router, handler := setupTestRouter()
	router.PUT("/pack-sizes/update", handler.UpdatePackSize)

	tests := []struct {
		name           string
		payload        map[string]int
		expectedStatus int
		expectedSizes  []int
	}{
		{
			name:           "Valid update",
			payload:        map[string]int{"oldSize": 250, "newSize": 300},
			expectedStatus: http.StatusOK,
			expectedSizes:  []int{300, 500, 1000, 2000, 5000},
		},
		{
			name:           "Non-existent old size",
			payload:        map[string]int{"oldSize": 999, "newSize": 300},
			expectedStatus: http.StatusBadRequest,
			expectedSizes:  nil,
		},
		{
			name:           "New size already exists",
			payload:        map[string]int{"oldSize": 250, "newSize": 500},
			expectedStatus: http.StatusBadRequest,
			expectedSizes:  nil,
		},
		{
			name:           "Invalid new size (zero)",
			payload:        map[string]int{"oldSize": 250, "newSize": 0},
			expectedStatus: http.StatusBadRequest,
			expectedSizes:  nil,
		},
		{
			name:           "Invalid new size (negative)",
			payload:        map[string]int{"oldSize": 250, "newSize": -100},
			expectedStatus: http.StatusBadRequest,
			expectedSizes:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/pack-sizes/update", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string][]int
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSizes, response["sizes"])
			}
		})
	}
}

func TestCalculatePacks(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/calculate", handler.CalculatePacks)

	tests := []struct {
		name           string
		quantity       string
		expectedStatus int
	}{
		{
			name:           "Valid quantity",
			quantity:       "500",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Zero quantity",
			quantity:       "0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Negative quantity",
			quantity:       "-100",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid quantity format",
			quantity:       "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/calculate?quantity="+tt.quantity, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "packs")
				assert.Contains(t, response, "totalItems")
			}
		})
	}
} 