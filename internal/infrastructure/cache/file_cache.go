package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache defines the cache storage interface.
type Cache interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, data []byte, ttl time.Duration) error
	Delete(key string) error
	Clear() error
}

type cacheEntry struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Data      []byte    `json:"data"`
}

// FileCache implements disk-based caching with TTL.
type FileCache struct {
	dir        string
	defaultTTL time.Duration
	mu         sync.RWMutex
}

// NewFileCache creates a new FileCache in the specified directory.
func NewFileCache(dir string, defaultTTL time.Duration) (*FileCache, error) {
	if dir == "" {
		dir = ".cache"
	}
	if defaultTTL <= 0 {
		defaultTTL = 1 * time.Hour
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory %q: %w", dir, err)
	}

	return &FileCache{
		dir:        dir,
		defaultTTL: defaultTTL,
	}, nil
}

func (fc *FileCache) keyPath(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:]) + ".json"
	return filepath.Join(fc.dir, filename)
}

// Get retrieves data for key if it exists and has not expired.
func (fc *FileCache) Get(key string) ([]byte, bool, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	path := fc.keyPath(key)
	fileData, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var entry cacheEntry
	if err := json.Unmarshal(fileData, &entry); err != nil {
		// Corrupted cache file, ignore and return miss
		_ = os.Remove(path)
		return nil, false, nil
	}

	if time.Now().After(entry.ExpiresAt) {
		// Expired
		go func(p string) {
			fc.mu.Lock()
			defer fc.mu.Unlock()
			_ = os.Remove(p)
		}(path)
		return nil, false, nil
	}

	return entry.Data, true, nil
}

// Set stores data for key with the given TTL.
func (fc *FileCache) Set(key string, data []byte, ttl time.Duration) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if ttl <= 0 {
		ttl = fc.defaultTTL
	}

	now := time.Now()
	entry := cacheEntry{
		Key:       key,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		Data:      data,
	}

	fileData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to serialize cache entry: %w", err)
	}

	path := fc.keyPath(key)
	if err := os.WriteFile(path, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write cache file %q: %w", path, err)
	}

	return nil
}

// Delete removes a specific cache entry.
func (fc *FileCache) Delete(key string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	path := fc.keyPath(key)
	return os.Remove(path)
}

// Clear removes all cached files in the cache directory.
func (fc *FileCache) Clear() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	entries, err := os.ReadDir(fc.dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			_ = os.Remove(filepath.Join(fc.dir, entry.Name()))
		}
	}

	return nil
}
