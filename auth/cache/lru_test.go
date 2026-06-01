package cache

import (
	"testing"
	"time"

	"github.com/Tentaculum-dev/go-sdk/auth"
	"github.com/google/uuid"
)

func id() *auth.Identity { return &auth.Identity{UserID: uuid.New()} }

func TestLRUGetSet(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", id(), time.Minute)
	if _, ok := c.Get("a"); !ok {
		t.Fatal("a should be present")
	}
	if _, ok := c.Get("missing"); ok {
		t.Fatal("missing should be absent")
	}
}

func TestLRUExpiry(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", id(), time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatal("a should have expired")
	}
}

func TestLRUEviction(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", id(), time.Minute)
	c.Set("b", id(), time.Minute)
	_, _ = c.Get("a")             // a now most-recently-used
	c.Set("c", id(), time.Minute) // evicts b (LRU)
	if _, ok := c.Get("b"); ok {
		t.Fatal("b should have been evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("a should survive")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("c should be present")
	}
}

func TestLRUZeroTTLNoop(t *testing.T) {
	c := NewLRU(2)
	c.Set("a", id(), 0)
	if _, ok := c.Get("a"); ok {
		t.Fatal("zero ttl should not store")
	}
}
