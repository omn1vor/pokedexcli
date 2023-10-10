package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "http://test1.org",
			val: []byte("raw data 1"),
		},
		{
			key: "http://test2.org",
			val: []byte("raw data 2"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case no. %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Fatal("expected to find key", c.key)
			}
			if string(val) != string(c.val) {
				t.Fatal("expected to find value equal to the one added")
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime * 3
	const key = "http://test1.org"

	cache := NewCache(baseTime)
	cache.Add(key, []byte("raw data1"))

	_, ok := cache.Get(key)
	if !ok {
		t.Fatal("expected to find key")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(key)
	if ok {
		t.Fatal("expected the cached value to be cleared")
	}
}
