// Example of an lru implementation using double linked lists (not thread safe)
package lrucache

import (
	"container/list"
	"errors"
	"strconv"
	"testing"
)

func mockCreator(key string) interface{} {
	return key
}

func doubleLinkedListContains(l *list.List, elm string) bool {
	for node := l.Front(); node != nil; node = node.Next() {
		if node.Value.(string) == elm {
			return true
		}
	}
	return false
}

func toString(a interface{}) string {
	if i, ok := a.(string); ok {
		return i
	}
	panic("cannot convert in string")
}

func TestAdd(t *testing.T) {
	lru := New(mockCreator, 10)
	if lru.add("key", "val") != "val" {
		t.Error("must returns the value passed")
	}

	if toString(lru.ddl.Front().Value) != "key" {
		t.Error("key must be ther first item after adding it")
	}
	if val, ok := lru.cache["key"]; !ok {
		t.Error("key not in cache")
	} else if toString(val.Value) != "val" {
		t.Error("bad value associated to key")
	}

}

func TestRemoveExtraKeys(t *testing.T) {
	lru := New(mockCreator, 4)
	for i := 0; i < 10; i++ {
		lru.add(strconv.Itoa(i), strconv.Itoa(i))
	}
	if lru.Size() != 10 {
		t.Error("bad initial length")
	}
	lru.removeExtraKeys()
	if lru.Size() != 4 {
		t.Error("bad final length.  should be 4 but have ", lru.Size())
	}
	if lru.ddl.Len() != 4 {
		t.Error("double linked list has incorrect length")
	}
}

func checkLRUCacheState(lru *LRUCache, expected ...string) error {
	expectedSize := len(expected)
	if len(lru.cache) != expectedSize {
		return errors.New("cache of wrong size")
	}

	if lru.ddl.Len() != expectedSize {
		return errors.New("double linked list of wrong size")
	}
	for _, e := range expected {
		if _, ok := lru.cache[e]; !ok {
			return errors.New("key " + e + " not in cache")
		}
		if !doubleLinkedListContains(lru.ddl, e) {
			return errors.New("key " + e + " not in double linked list")
		}
	}
	return nil
}

func TestGet(t *testing.T) {
	lru := New(mockCreator, 2)
	if v := lru.Get("0"); v != "0" {
		t.Error()
	}
	if err := checkLRUCacheState(lru, "0"); err != nil {
		t.Error(err)
		return
	}

	lru.Get("1")
	if err := checkLRUCacheState(lru, "1", "0"); err != nil {
		t.Error(err)
		return
	}

	lru.Get("2")
	if err := checkLRUCacheState(lru, "2", "1"); err != nil {
		t.Error(err)
		return
	}

	lru.Get("1")
	if err := checkLRUCacheState(lru, "1", "2"); err != nil {
		t.Error(err)
		return
	}

	lru.Get("3")
	if err := checkLRUCacheState(lru, "3", "1"); err != nil {
		t.Error(err)
		return
	}
}
