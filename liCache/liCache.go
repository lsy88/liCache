package liCache

import (
	"github.com/lsy88/liCache/cache"
	"time"
)

type cacheServer struct {
	memCache cache.Cache
}

func NewMemCache() *cacheServer {
	return &cacheServer{
		memCache: cache.NewMemCache(),
	}
}

func (cs *cacheServer) SetMaxMemory(size string) bool {
	return cs.memCache.SetMaxMemory(size)
}

func (cs *cacheServer) Set(key string, val interface{}, expire ...time.Duration) bool {
	expireTs := time.Second * 0
	if len(expire) > 0 {
		expireTs = expire[0]
	}
	return cs.memCache.Set(key, val, expireTs)
}

func (cs *cacheServer) Get(key string) (interface{}, bool) {
	return cs.memCache.Get(key)
}

func (cs *cacheServer) Del(key string) {
	cs.memCache.Del(key)
}

func (cs *cacheServer) Exist(key string) bool {
	return cs.memCache.Exist(key)
}

func (cs *cacheServer) Flush() bool {
	return cs.memCache.Flush()
}

func (cs *cacheServer) Keys() int64 {
	return cs.memCache.Keys()
}
