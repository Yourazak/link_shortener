package storage

import (
	"testing"
)

func TestCache_Set_Success(t *testing.T) {
	cache := NewCache()
	err := cache.Set("short1", "https://example.com")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	val, exists := cache.DB["short1"]
	if !exists || val != "https://example.com" {
		t.Errorf("Expected 'https://example.com', got '%s'", val)
	}
}

func TestCache_Set_AlreadyExists(t *testing.T) {
	cache := NewCache()
	_ = cache.Set("short1", "https://example.com")

	err := cache.Set("short1", "https://new-url.com")
	if err == nil || err.Error() != "key is already exists" {
		t.Errorf("Expected error 'key is already exists', got %v", err)
	}
}

func TestCache_Get_Success(t *testing.T) {
	cache := NewCache()
	_ = cache.Set("short1", "https://example.com")

	value, err := cache.Get("short1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value != "https://example.com" {
		t.Errorf("Expected 'https://example.com', got '%s'", value)
	}
}

func TestCache_Get_NotFound(t *testing.T) {
	cache := NewCache()

	_, err := cache.Get("notfound")
	if err == nil || err.Error() != "key doest exists" {
		t.Errorf("Expected error 'key doest exists', got %v", err)
	}
}

func TestCache_Close(t *testing.T) {
	cache := NewCache()
	err := cache.Close()
	if err != nil {
		t.Errorf("Close() should return nil, got %v", err)
	}
}
