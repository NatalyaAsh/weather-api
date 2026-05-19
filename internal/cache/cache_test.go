package cache

import (
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	cache := New(1 * time.Second)

	cache.Set("key1", "value1")
	cache.Set("key2", 42)

	val, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}

	val, found = cache.Get("key2")
	if !found {
		t.Error("Expected to find key2")
	}
	if val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}
}

func TestCacheNotFound(t *testing.T) {
	cache := New(1 * time.Second)

	_, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
}

func TestCacheExpiration(t *testing.T) {
	cache := New(100 * time.Millisecond)

	cache.Set("key", "value")

	// Сразу после установки должно быть
	_, found := cache.Get("key")
	if !found {
		t.Error("Key should exist immediately after set")
	}

	// Ждём истечения
	time.Sleep(150 * time.Millisecond)

	_, found = cache.Get("key")
	if found {
		t.Error("Key should have expired")
	}
}

func TestCacheOverwrite(t *testing.T) {
	cache := New(1 * time.Second)

	cache.Set("key", "first")
	cache.Set("key", "second")

	val, found := cache.Get("key")
	if !found {
		t.Error("Expected to find key")
	}
	if val != "second" {
		t.Errorf("Expected 'second', got %v", val)
	}
}
