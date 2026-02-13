package service

import (
	"errors"
	"testing"

	"ozon_link_shortener/internal/tests/mocks"

	"github.com/golang/mock/gomock"
)

func TestShortener(t *testing.T) {
	url := "https://example.com"
	short := Shortener(url)

	if len(short) != 10 {
		t.Errorf("Expected length 10, got %d", len(short))
	}
}

func TestWriteToStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockStorage(ctrl)
	mockDB := mocks.NewMockStorage(ctrl)

	service := &Service{Cache: mockCache, DB: mockDB, Conf: &Config{Storage: "postgres"}}

	shortUrl := "abc123"
	req := Request{LongUrl: "https://example.com"}

	mockCache.EXPECT().Set(shortUrl, req.LongUrl).Return(nil)
	mockDB.EXPECT().Set(shortUrl, req.LongUrl).Return(nil)

	err := WriteToStorage(service, shortUrl, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestGetFromStorage_CacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockStorage(ctrl)
	mockDB := mocks.NewMockStorage(ctrl)
	service := &Service{Cache: mockCache, DB: mockDB, Conf: &Config{Storage: "postgres"}}

	shortUrl := "abc123"
	longUrl := "https://example.com"

	mockCache.EXPECT().Get(shortUrl).Return(longUrl, nil)

	result, err := GetFromStorage(service, shortUrl)
	if err != nil || result != longUrl {
		t.Errorf("Expected %s, got %s (err: %v)", longUrl, result, err)
	}
}

func TestGetFromStorage_DBHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockStorage(ctrl)
	mockDB := mocks.NewMockStorage(ctrl)
	service := &Service{Cache: mockCache, DB: mockDB, Conf: &Config{Storage: "postgres"}}

	shortUrl := "abc123"
	longUrl := "https://example.com"

	mockCache.EXPECT().Get(shortUrl).Return("", errors.New("cache miss"))
	mockDB.EXPECT().Get(shortUrl).Return(longUrl, nil)

	result, err := GetFromStorage(service, shortUrl)
	if err != nil || result != longUrl {
		t.Errorf("Expected %s, got %s (err: %v)", longUrl, result, err)
	}
}

func TestGetFromStorage_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockStorage(ctrl)
	mockDB := mocks.NewMockStorage(ctrl)
	service := &Service{Cache: mockCache, DB: mockDB, Conf: &Config{Storage: "postgres"}}

	shortUrl := "abc123"

	mockCache.EXPECT().Get(shortUrl).Return("", errors.New("cache miss"))
	mockDB.EXPECT().Get(shortUrl).Return("", errors.New("db miss"))

	_, err := GetFromStorage(service, shortUrl)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
