package tests

import (
	"fmt"
	"testing"
	"time"

	pokecache "github.com/KrupaH/golang-pokedex/internal/pokecache"
)

func TestAddGetReaploop(t *testing.T) {
	const interval = 1 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
			}
			if string(val) != string(c.val) {
				t.Errorf("expected value not found")
			}
			time.Sleep(interval + 1)
			val, ok = cache.Get(c.key)
			if ok {
				t.Errorf("Expected key %v to be cleared by now", c.key)
			}
		})
	}
}
