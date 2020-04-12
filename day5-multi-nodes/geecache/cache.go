package geecache

import (
	lru2 "awesomeProject4/GeeCaches/day2-single-node/geecache/lru"
	"sync"
)

type cache struct {
	mu sync.Mutex
	lru *lru2.Cache
	cacheBytes int64
}

//支持并发的读写
func(c *cache) add(key string, value ByteView){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		c.lru = lru2.New(c.cacheBytes,nil)//延迟初始化，在使用时才创建
	}
	c.lru.Add(key,value)
}

func(c *cache) get(key string)(value ByteView,ok bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil{
		return
	}
	if v,ok := c.lru.Get(key);ok{
		return v.(ByteView),ok
	}
	return

}
