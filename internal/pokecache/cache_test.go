package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	c := NewCache(50 * time.Millisecond)

	key := "https://pokeapi.co/test"
	val := []byte("hello")

	c.Add(key, val)

	got, ok := c.Get(key)
	if !ok {
		t.Fatalf("expected key to exist")
	}

	if string(got) != "hello" {
		t.Fatalf("expected %q, got %q", "hello", string(got))
	}
}

func TestCacheReap(t *testing.T) {
	c := NewCache(10 * time.Millisecond)

	key := "old"
	val := []byte("data")

	c.Add(key, val)

	time.Sleep(30 * time.Millisecond)

	_, ok := c.Get(key)
	if ok {
		t.Fatalf("expected entry to be reaped")
	}
}
