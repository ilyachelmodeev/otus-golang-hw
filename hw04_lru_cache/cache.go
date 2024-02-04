package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       sync.Mutex{},
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	currentItem, ok := c.items[key]
	if ok {
		currentItem.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(currentItem)
		return true
	}

	if len(c.items) == c.capacity {
		tail := c.queue.Back()
		keyToDel := tail.Value.(cacheItem).key
		c.queue.Remove(tail)
		delete(c.items, keyToDel)
	}

	cacheItemValue := cacheItem{key: key, value: value}
	newCurrentItem := c.queue.PushFront(cacheItemValue)
	c.items[key] = newCurrentItem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if currentItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(currentItem)
		return currentItem.Value.(cacheItem).value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
