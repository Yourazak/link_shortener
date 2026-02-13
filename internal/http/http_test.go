package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"ozon_link_shortener/internal/service"
	"ozon_link_shortener/internal/tests/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSaveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheMock := mocks.NewMockStorage(ctrl)
	DBMock := mocks.NewMockStorage(ctrl)

	cacheMock.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	DBMock.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockService := &service.Service{Cache: cacheMock, DB: DBMock, Conf: &service.Config{Storage: "postgres"}}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/", func(ctx *gin.Context) {
		SaveUrl(ctx, mockService)
	})

	reqBody := map[string]string{"url": "https://example.com"}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), `"url":`)
}

func TestGetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorage(ctrl)
	mockService := &service.Service{Cache: mockStorage, DB: mockStorage}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/:shortUrl", func(ctx *gin.Context) {
		GetUrl(ctx, mockService)
	})

	mockStorage.EXPECT().Get("short123").Return("https://example.com", nil)

	req := httptest.NewRequest(http.MethodGet, "/short123", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusMovedPermanently, resp.Code)

	assert.Equal(t, "https://example.com", resp.Header().Get("Location"))
}

func TestGetUrl_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheMock := mocks.NewMockStorage(ctrl)
	DBMock := mocks.NewMockStorage(ctrl)

	cacheMock.EXPECT().Get("notfound").Return("", errors.New("key does not exists"))
	DBMock.EXPECT().Get("notfound").Return("", errors.New("key does not exists"))

	mockService := &service.Service{Cache: cacheMock, DB: DBMock, Conf: &service.Config{Storage: "postgres"}}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/:shortUrl", func(ctx *gin.Context) {
		GetUrl(ctx, mockService)
	})

	req := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
