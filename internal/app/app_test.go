package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bidntb/metrics/internal/middleware"
	"bidntb/metrics/internal/service/handler"
	"bidntb/metrics/internal/service/metrics"
	"bidntb/metrics/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	storageInstance := storage.NewMemStorage()
	metricsSvc := metrics.NewService(storageInstance)
	h := handler.NewHandler(metricsSvc)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler(h.NotFoundHandler, h.BadRequestHandler))

	r.GET("/", h.ListMetrics)
	r.POST("/update/", h.UpdateMetric)
	r.POST("/value/", h.GetMetricJSON)
	r.POST("/update/gauge/:name/:value", middleware.ValidateGaugeValue(), h.UpdateGauge)
	r.POST("/update/counter/:name/:value", middleware.ValidateCounterValue(), h.UpdateCounter)
	r.GET("/value/:type/:name", h.GetValue)

	r.POST("/update/counter", h.NotFoundHandler)
	r.POST("/update/gauge/", h.NotFoundHandler)
	r.POST("/update/gauge", h.NotFoundHandler)
	r.POST("/update/:wrong", h.BadRequestHandler)
	r.POST("/update/:wrong/*any", h.BadRequestHandler)
	r.NoRoute(h.NotFoundHandler)

	return r
}

func TestServerRoutes(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name         string
		method       string
		path         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "List all metrics",
			method:       http.MethodGet,
			path:         "/",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Update gauge JSON",
			method:       http.MethodPost,
			path:         "/update/",
			body:         `{"type":"gauge","name":"test_gauge","value":123.45}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Update counter path params",
			method:       http.MethodPost,
			path:         "/update/counter/test_counter/5",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get metric value path params",
			method:       http.MethodGet,
			path:         "/value/gauge/test_gauge",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid value - bad request",
			method:       http.MethodPost,
			path:         "/update/counter/test/none",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Not found route",
			method:       http.MethodPost,
			path:         "/unknown",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestUpdateMetricJSON(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"type":"gauge","name":"cpu_load","value":189736.689}`
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp metrics.MetricResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.Equal(t, "gauge", resp.MType)
	assert.NotEmpty(t, resp.ID)
	assert.InDelta(t, 189736.689, resp.Value, 0.001)
}

func TestCounterAccumulation(t *testing.T) {
	r := setupTestRouter()

	reqs := []string{
		`{"type":"counter","name":"requests","value":10}`,
		`{"type":"counter","name":"requests","value":20}`,
	}

	for _, body := range reqs {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/value/counter/requests", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "30", w.Body.String())
}

func TestGetMetricByID(t *testing.T) {
	r := setupTestRouter()

	w := httptest.NewRecorder()
	reqBody := `{"type":"gauge","name":"test_id","value":42.0}`
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var createResp metrics.MetricResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &createResp))
	metricID := createResp.ID

	w = httptest.NewRecorder()
	getBody := fmt.Sprintf(`{"id":"%s","MType":"gauge"}`, metricID)
	req, _ = http.NewRequest(http.MethodPost, "/value/", bytes.NewBufferString(getBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResp metrics.MetricResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &getResp))
	assert.Equal(t, metricID, getResp.ID)
}

func TestValidationErrors(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name         string
		path         string
		body         string
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Missing JSON",
			path:         "/update/",
			body:         `{"type":"gauge"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid counter value",
			path:         "/update/counter/test/none",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Counter value must be integer",
		},
		{
			name:         "Invalid gauge value",
			path:         "/update/gauge/test/none",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Gauge value must be number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, tt.path, bytes.NewBufferString(tt.body))
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedMsg != "" {
				assert.Contains(t, w.Body.String(), tt.expectedMsg)
			}
		})
	}
}
