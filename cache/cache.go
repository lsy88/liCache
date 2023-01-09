package cache

import "time"

type Cache interface {
	SetMaxMemory(string) bool
	Set(key string, val interface{}, expire time.Duration) bool
	Get(key string) (interface{}, bool)
	Del(key string)
	Exist(key string) bool
	Flush() bool //清空所有的key
	Keys() int64 //获取缓存中所有key的数量
}
