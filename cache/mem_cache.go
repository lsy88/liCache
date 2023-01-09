package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type memCache struct {
	maxMemorySize int64
	
	maxMemorySizeStr string //最大内存字符串表示
	
	usedMemorySize int64 //当前已使用内存
	
	B map[string]*memCacheValue //缓存结构
	
	mutex sync.RWMutex //读写锁，允许并发
	
	clearExpireItemTimeInterval time.Duration //清除过期缓存时间间隔
}

type memCacheValue struct {
	val interface{}
	//过期时间
	expireTime time.Time
	//有效时长
	keep time.Duration
	//size 每个value的大小
	size int64
}

func NewMemCache() Cache {
	mc := &memCache{
		B:                           make(map[string]*memCacheValue),
		clearExpireItemTimeInterval: time.Second * 2,
	}
	go mc.clearExpireItem()
	return mc
}

func (mc *memCache) SetMaxMemory(size string) bool {
	maxMemorySize, maxMemorySizeStr := parseSize(size)
	mc.maxMemorySize = maxMemorySize
	mc.maxMemorySizeStr = maxMemorySizeStr
	return true
}

func (mc *memCache) Set(key string, val interface{}, expire time.Duration) bool {
	//写锁控制
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	v := &memCacheValue{
		val:        val,
		expireTime: time.Now().Add(expire),
		keep:       expire,
		size:       getValSize(val),
	}
	mc.Del(key)
	mc.add(key, v)
	if mc.usedMemorySize > mc.maxMemorySize {
		mc.Del(key)
		//todo 可以选择不panic,删除掉一些过期的key
		log.Println(fmt.Sprintf("max memory size %s", mc.maxMemorySizeStr))
		return false
	}
	return true
}

func (mc *memCache) get(key string) (*memCacheValue, bool) {
	val, ok := mc.B[key]
	return val, ok
}

func (mc *memCache) Get(key string) (interface{}, bool) {
	//读锁控制
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	get, ok := mc.get(key)
	if ok {
		//判断是否过期
		if get.keep != 0 && get.expireTime.Before(time.Now()) {
			mc.del(key)
			return nil, false
		}
		return get.val, ok
	}
	return nil, false
}

func (mc *memCache) del(key string) {
	value, ok := mc.get(key)
	if ok && value != nil {
		mc.usedMemorySize -= value.size
		delete(mc.B, key)
	}
}

func (mc *memCache) Del(key string) {
	//写锁控制
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.del(key)
}

func (mc *memCache) add(key string, val *memCacheValue) {
	mc.B[key] = val
	mc.usedMemorySize += val.size
}

func (mc *memCache) Exist(key string) bool {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	
	_, ok := mc.get(key)
	return ok
}

//清空所有的key
func (mc *memCache) Flush() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	mc.B = make(map[string]*memCacheValue, 0)
	mc.usedMemorySize = 0
	return true
}

//获取缓存中所有key的数量
func (mc *memCache) Keys() int64 {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	return int64(len(mc.B))
}

func (mc *memCache) clearExpireItem() {
	ticker := time.NewTicker(mc.clearExpireItemTimeInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			for key, item := range mc.B {
				if item.keep != 0 && time.Now().After(item.expireTime) {
					mc.mutex.Lock()
					mc.del(key)
					mc.mutex.Unlock()
				}
			}
		default:
		
		}
	}
}
