package lru

import "container/list"
//这个版本不支持并发
type Cache struct {
	maxBytes int64//最大允许内存
	nbytes int64//已经使用内存
	ll *list.List//双向链表
	cache map[string]*list.Element
	OnEvicted func(key string,value Value)//某条记录被移出时的回调函数
}

type entry struct {
	//双向链表节点的数据类型
	key string
	value Value
}
type Value interface {
	Len() int
	//值所占用的内存大小
}
func New(maxBytes int64,onEvicted func(string,Value)) *Cache{
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
//查找功能
func(c *Cache) Get(key string)(value Value,ok bool){
	if ele,ok := c.cache[key];ok {
		c.ll.MoveToFront(ele)//约定front为队尾
		kv := ele.Value.(*entry)//??
		return kv.value,true
	}
	return
}
//删除功能
func(c *Cache) RemoveOldest(){
	ele := c.ll.Back()
	if ele != nil{
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache,kv.key)
		c.nbytes -= int64(len(kv.key))+int64(kv.value.Len())
		if c.OnEvicted != nil{
			c.OnEvicted(kv.key,kv.value)
		}
	}
}
//新增，修改功能
func(c *Cache) Add(key string,value Value){
	if ele,ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len())- int64(kv.value.Len())//可能更新值
		kv.value = value
	} else{
		ele := c.ll.PushFront(&entry{key,value})//?
		c.cache[key] = ele
		c.nbytes += int64(len(key))+int64(value.Len())
	}
	for c.maxBytes!=0 && c.maxBytes<c.nbytes{
		c.RemoveOldest()
	}
}
//缓存数量
func(c *Cache) Len() int{
	return c.ll.Len()
}
