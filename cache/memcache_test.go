package cache

import (
	"testing"
	"time"
)

func TestCacheOP(t *testing.T) {
	data := []struct {
		key    string
		val    interface{}
		expire time.Duration
	}{
		{"zhangsan", 678, 10},
		{"lisi", 238, 15},
		{"wangsa", 8, 0},
		{"kkxing", 6, 20},
		{"sssca", false, 10},
		{"wwher", "ssssaa", 15},
		{"w3asc", []int{1, 2, 3}, 0},
		{"httrjr", "wafav", 20},
	}
	c := NewMemCache()
	c.SetMaxMemory("10MB")
	for _, item := range data {
		c.Set(item.key, item.val, item.expire)
		get, ok := c.Get(item.key)
		if !ok {
			t.Error("error")
		}
		t.Log(get)
	}
}
