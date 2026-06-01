// Package cache provides an in-memory LRU implementation of auth.TokenCache.
package cache

import (
	"container/list"
	"sync"
	"time"

	"github.com/ViitoJooj/sdk/auth"
)

type entry struct {
	key       string
	id        *auth.Identity
	expiresAt time.Time
}

// LRU is a fixed-capacity, TTL-aware token cache. Safe for concurrent use.
type LRU struct {
	mu    sync.Mutex
	cap   int
	ll    *list.List               // front = most recently used
	items map[string]*list.Element // key -> element holding *entry
}

// NewLRU returns an LRU holding up to capacity entries. capacity <= 0 panics.
func NewLRU(capacity int) *LRU {
	if capacity <= 0 {
		panic("cache: capacity must be > 0")
	}
	return &LRU{
		cap:   capacity,
		ll:    list.New(),
		items: make(map[string]*list.Element, capacity),
	}
}

// Get returns the cached identity if present and not expired.
func (c *LRU) Get(key string) (*auth.Identity, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return nil, false
	}
	en := el.Value.(*entry)
	if time.Now().After(en.expiresAt) {
		c.removeElement(el)
		return nil, false
	}
	c.ll.MoveToFront(el)
	return en.id, true
}

// Set stores id under key with the given TTL, evicting the least recently
// used entry when over capacity. Non-positive TTL is a no-op.
func (c *LRU) Set(key string, id *auth.Identity, ttl time.Duration) {
	if ttl <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	exp := time.Now().Add(ttl)
	if el, ok := c.items[key]; ok {
		en := el.Value.(*entry)
		en.id = id
		en.expiresAt = exp
		c.ll.MoveToFront(el)
		return
	}
	el := c.ll.PushFront(&entry{key: key, id: id, expiresAt: exp})
	c.items[key] = el
	if c.ll.Len() > c.cap {
		if back := c.ll.Back(); back != nil {
			c.removeElement(back)
		}
	}
}

// Len returns the current number of entries (including any expired but
// not-yet-evicted ones).
func (c *LRU) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

func (c *LRU) removeElement(el *list.Element) {
	c.ll.Remove(el)
	delete(c.items, el.Value.(*entry).key)
}

// compile-time check that *LRU satisfies auth.TokenCache.
var _ auth.TokenCache = (*LRU)(nil)
