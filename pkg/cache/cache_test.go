package cache

import (
	"net/http"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	cache := New()

	// Test case 1: Get an item that does not exist
	_, found := cache.Get("nonexistent")
	if found {
		t.Errorf("expected item to not be found")
	}

	// Test case 2: Get an item that exists and is not expired
	item := &Item{
		ResponseBody:       []byte("response body"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(1 * time.Hour),
	}
	cache.Set("key1", item)

	retrievedItem, found := cache.Get("key1")
	if !found {
		t.Errorf("expected item to be found")
	}
	if string(retrievedItem.ResponseBody) != "response body" {
		t.Errorf("expected response body to be 'response body', got '%s'", string(retrievedItem.ResponseBody))
	}
	if retrievedItem.ResponseStatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, retrievedItem.ResponseStatusCode)
	}

	// Test case 3: Get an item that exists but is expired
	expiredItem := &Item{
		ResponseBody:       []byte("expired response body"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(-1 * time.Hour),
	}
	cache.Set("key2", expiredItem)

	_, found = cache.Get("key2")
	if found {
		t.Errorf("expected expired item to not be found")
	}
}
func TestCache_Set(t *testing.T) {
	cache := New()

	// Test case 1: Set an item and retrieve it
	item := &Item{
		ResponseBody:       []byte("response body"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(1 * time.Hour),
	}
	cache.Set("key1", item)

	retrievedItem, found := cache.Get("key1")
	if !found {
		t.Errorf("expected item to be found")
	}
	if string(retrievedItem.ResponseBody) != "response body" {
		t.Errorf("expected response body to be 'response body', got '%s'", string(retrievedItem.ResponseBody))
	}
	if retrievedItem.ResponseStatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, retrievedItem.ResponseStatusCode)
	}

	// Test case 2: Set an item when the cache is full
	cache.capacity = 1
	item2 := &Item{
		ResponseBody:       []byte("new response body"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(1 * time.Hour),
	}
	cache.Set("key2", item2)

	_, found = cache.Get("key1")
	if found {
		t.Errorf("expected the first item to be evicted")
	}

	retrievedItem, found = cache.Get("key2")
	if !found {
		t.Errorf("expected the new item to be found")
	}
	if string(retrievedItem.ResponseBody) != "new response body" {
		t.Errorf("expected response body to be 'new response body', got '%s'", string(retrievedItem.ResponseBody))
	}
	if retrievedItem.ResponseStatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, retrievedItem.ResponseStatusCode)
	}
}
func TestCache_RemoveAll(t *testing.T) {
	cache := New()

	// Add some items to the cache
	item1 := &Item{
		ResponseBody:       []byte("response body 1"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(1 * time.Hour),
	}
	cache.Set("key1", item1)

	item2 := &Item{
		ResponseBody:       []byte("response body 2"),
		ResponseHeaders:    http.Header{"Content-Type": []string{"application/json"}},
		ResponseStatusCode: http.StatusOK,
		Expiration:         time.Now().Add(1 * time.Hour),
	}
	cache.Set("key2", item2)

	// Ensure items are in the cache
	if _, found := cache.Get("key1"); !found {
		t.Errorf("expected item1 to be found")
	}
	if _, found := cache.Get("key2"); !found {
		t.Errorf("expected item2 to be found")
	}

	// Remove all items from the cache
	cache.RemoveAll()

	// Ensure cache is empty
	if _, found := cache.Get("key1"); found {
		t.Errorf("expected item1 to be removed")
	}
	if _, found := cache.Get("key2"); found {
		t.Errorf("expected item2 to be removed")
	}
	if len(cache.itemsMap) != 0 {
		t.Errorf("expected itemsMap to be empty, got %d items", len(cache.itemsMap))
	}
	if cache.itemsList.Len() != 0 {
		t.Errorf("expected itemsList to be empty, got %d items", cache.itemsList.Len())
	}
}


