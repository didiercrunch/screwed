package lrucache

import (
	"container/list"
)

type Creator func(key string) interface{}

type cacheKey struct {
	Value                interface{}
	DoubleLinkedListNode *list.Element
}

type LRUCache struct {
	cache   map[string]*cacheKey
	ddl     *list.List
	creator Creator
	length  int
}

func New(creator Creator, length int) *LRUCache {
	return &LRUCache{
		cache:   make(map[string]*cacheKey),
		creator: creator,
		ddl:     new(list.List).Init(),
		length:  length,
	}
}

func (lru *LRUCache) add(key string, val interface{}) interface{} {
	lru.ddl.PushFront(key)
	lru.cache[key] = &cacheKey{val, lru.ddl.Front()}
	return val
}

func (lru *LRUCache) Size() int {
	return len(lru.cache)
}

func (lru *LRUCache) removeExtraKeys() {
	if lru.ddl.Len() < lru.length {
		return
	}
	for e := lru.ddl.Back(); lru.ddl.Len() > lru.length; {
		delete(lru.cache, e.Value.(string))
		be := e.Prev()
		lru.ddl.Remove(e)
		e = be
	}
}

func (lru *LRUCache) Get(key string) interface{} {
	if val, ok := lru.cache[key]; ok {
		lru.ddl.MoveToFront(val.DoubleLinkedListNode)
		return val.Value
	}
	ret := lru.add(key, lru.creator(key))
	lru.removeExtraKeys()
	return ret
}
