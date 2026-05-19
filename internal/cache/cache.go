package cache

import (
	"log/slog"
	"sync"
	"time"
)

// Item хранит значение и время истечения
type Item struct {
	Value     interface{} // здесь будем хранить готовую WeatherResponse
	ExpiresAt time.Time
}

// Cache основная структура кеша
type Cache struct {
	items map[string]Item
	mu    sync.RWMutex
	ttl   time.Duration // время жизни кеша
}

// New создаёт новый кеш с заданным TTL (например, 5 минут)
func New(ttl time.Duration) *Cache {
	return &Cache{
		items: make(map[string]Item),
		ttl:   ttl,
	}
}

// Set сохраняет значение по ключу
func (c *Cache) Set(key string, value interface{}) {
	slog.Debug("Cache set", "key", key)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get возвращает значение и true, если ключ существует и не истёк
func (c *Cache) Get(key string) (interface{}, bool) {
	slog.Debug("Cache get", "key", key)
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// проверяем, не истекло ли
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Value, true
}
