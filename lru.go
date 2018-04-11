package lru

import "github.com/gotoolkit/list"

// Cache LRU缓存，线程不安全
type Cache struct {
	// MaxEntries 缓存最大条目
	MaxEntries int
	// OnEvicted 清楚条目时执行
	OnEvicted func(key Key, value interface{})

	ll    *list.List
	cache map[interface{}]*list.Node
}

// Key 可比较的类型
type Key interface{}

type entry struct {
	key   Key
	value interface{}
}

// New 创建缓存
func New(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Node),
	}
}

// Add 缓存添加键值
func (c *Cache) Add(key Key, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Node)
		c.ll = list.New()
	}
	if n, ok := c.cache[key]; ok {
		c.ll.MoveToFront(n)
		n.Value.(*entry).value = value
		return
	}
	nln := c.ll.PushFront(&entry{key, value})
	c.cache[key] = nln
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get 缓存获取键值
func (c *Cache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if nln, hit := c.cache[key]; hit {
		c.ll.MoveToFront(nln)
		return nln.Value.(*entry).value, true
	}
	return
}

// Remove 缓存移除键值
func (c *Cache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if nln, hit := c.cache[key]; hit {
		c.removeNode(nln)
	}
}

// Len 缓存长度
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear 缓存清除
func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, n := range c.cache {
			kv := n.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}

// RemoveOldest 缓存移除最老的节点
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	nln := c.ll.Back()
	if nln != nil {
		c.removeNode(nln)
	}
}

func (c *Cache) removeNode(n *list.Node) {
	c.ll.Remove(n)
	kv := n.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}
