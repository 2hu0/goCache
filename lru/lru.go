package lru

import (
	"container/list"
)

// Cache is a LRU cache. It is not safe for concurrent access.

type Cache struct {
	//允许使用的最大内存
	maxBytes int64
	//当前已经使用的内存
	nbytes int64
	ll     *list.List
	//key是字符串 value是对应链表中节点的指针
	cache map[string]*list.Element
	// 某条记录被移除时的回调函数
	OnEvicted func(key string, value Value)
}

// 双向链表节点的数据类型 ，保存key的好处：淘汰节点时候，在字典中删除对应的映射
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// cache的构造器
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找功能： 1:从字典中找到对应的双向链表节点 2:将该节点移动到尾部,front是尾部
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除:移除队首 (最近访问最少)
func (c *Cache) RemoveOldest() {
	//获取队首节点
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//删除映射关系
		delete(c.cache, kv.key)
		//更新内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增，修改功能
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
