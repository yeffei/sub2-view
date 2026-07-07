//go:build embed

package web

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// HTMLCache manages the cached index.html with injected settings
type HTMLCache struct {
	mu              sync.RWMutex
	cachedHTML      map[string][]byte
	etags           map[string]string
	baseHTMLHash    string // Hash of the original index.html (immutable after build)
	settingsVersion uint64 // Incremented when settings change
}

// CachedHTML represents the cache state
type CachedHTML struct {
	Content []byte
	ETag    string
}

// NewHTMLCache creates a new HTML cache instance
func NewHTMLCache() *HTMLCache {
	return &HTMLCache{
		cachedHTML: make(map[string][]byte),
		etags:      make(map[string]string),
	}
}

// SetBaseHTML initializes the cache with the base HTML template
func (c *HTMLCache) SetBaseHTML(baseHTML []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	hash := sha256.Sum256(baseHTML)
	c.baseHTMLHash = hex.EncodeToString(hash[:8]) // First 8 bytes for brevity
}

// Invalidate marks the cache as stale
func (c *HTMLCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.settingsVersion++
	c.cachedHTML = make(map[string][]byte)
	c.etags = make(map[string]string)
}

// Get returns the cached HTML or nil if cache is stale
func (c *HTMLCache) Get(keys ...string) *CachedHTML {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := normalizeCacheKey(keys...)
	if c.cachedHTML == nil {
		return nil
	}
	html := c.cachedHTML[key]
	etag := c.etags[key]
	if html == nil || etag == "" {
		return nil
	}
	return &CachedHTML{Content: html, ETag: etag}
}

// Set updates the cache with new rendered HTML
func (c *HTMLCache) Set(html []byte, settingsJSON []byte, keys ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cachedHTML == nil {
		c.cachedHTML = make(map[string][]byte)
	}
	if c.etags == nil {
		c.etags = make(map[string]string)
	}
	key := normalizeCacheKey(keys...)
	c.cachedHTML[key] = html
	c.etags[key] = c.generateETag(settingsJSON, key)
}

// generateETag creates an ETag from base HTML hash + settings hash
func (c *HTMLCache) generateETag(settingsJSON []byte, key string) string {
	settingsHash := sha256.Sum256(settingsJSON)
	keyHash := sha256.Sum256([]byte(key))
	return `"` + c.baseHTMLHash + "-" + hex.EncodeToString(settingsHash[:6]) + "-" + hex.EncodeToString(keyHash[:6]) + `"`
}

func normalizeCacheKey(keys ...string) string {
	if len(keys) == 0 || keys[0] == "" {
		return "default"
	}
	return keys[0]
}
