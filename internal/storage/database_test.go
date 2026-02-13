package storage

import (
	"context"
	"testing"

	"ozon_link_shortener/internal/tests/mocks"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
)

func TestNewDB(t *testing.T) {
	connStr := "postgres://postgres:postgres@localhost:5433/postgres"
	db, err := NewDB(connStr)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if db == nil {
		t.Errorf("Database connection is nil")
	}

	if db != nil {
		db.Close()
	}
}

func TestDatabase_Get_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPgxConnInterface(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	db := &Database{DB: mockDB}

	mockRow.EXPECT().Scan(gomock.Any()).SetArg(0, "https://example.com").Return(nil)

	mockDB.EXPECT().
		QueryRow(context.Background(), `SELECT value FROM url_table WHERE key=$1`, "short1").
		Return(mockRow)

	result, err := db.Get("short1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "https://example.com" {
		t.Errorf("Expected 'https://example.com', got '%s'", result)
	}
}

func TestDatabase_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPgxConnInterface(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	db := &Database{DB: mockDB}

	mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows)

	mockDB.EXPECT().
		QueryRow(context.Background(), `SELECT value FROM url_table WHERE key=$1`, "notfound").
		Return(mockRow)

	_, err := db.Get("notfound")
	if err == nil || err.Error() != "Key dosnt exists in database" {
		t.Errorf("Expected 'Key dosnt exists in database', got %v", err)
	}
}

func TestDatabase_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockPgxConnInterface(ctrl)
	db := &Database{DB: mockDB}

	mockDB.EXPECT().Close(context.Background()).Return(nil)

	err := db.Close()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
