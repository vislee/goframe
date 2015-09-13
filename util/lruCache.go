// Copyright 2015 The liwq. All rights reserved.
package util

import (
	"container/list"
	"errors"
	"fmt"
	"time"
)

type Item struct {
	Value  interface{}
	expire int64
}

type LruCache struct {
	size   int
	expire int64
	l      *list.List
	cache  map[string]Item
}

func NewLruCache(l int, exp int64) *LruCache {
	if exp == 0 {
		exp = 30 * 24 * 3600
	}
	return &LruCache{size: l, expire: exp}
}

func (lc *LruCache) Init() {
	lc.l = list.New()
	lc.cache = make(map[string]Item, lc.size+1)
}

func (lc *LruCache) Update(key string, v interface{}) {
	_, ok := lc.cache[key]
	if ok {
		for lt := lc.l.Front(); lt != nil; lt = lt.Next() {
			if lt.Value == key {
				lc.l.MoveToFront(lt)
				break
			}
		}
		return
	}

	ntstmp := time.Now().Unix()
	lc.cache[key] = Item{Value: v, expire: ntstmp + lc.expire}
	if lc.l.Len() == 0 {
		lc.l.PushFront(key)
	} else {
		lc.l.InsertAfter(key, lc.l.Front())
	}

	if lc.l.Len() == lc.size+1 {
		li := lc.l.Back()
		lc.l.Remove(li)
		delete(lc.cache, key)
	}

	return
}

func (lc *LruCache) Search(key string) (interface{}, error) {
	itm, ok := lc.cache[key]
	if !ok {
		return nil, errors.New("cache miss")
	}

	ntsmp := time.Now().Unix()
	if itm.expire < ntsmp {
		for li := lc.l.Front(); li != nil; li = li.Next() {
			if li.Value == key {
				lc.l.Remove(li)
				delete(lc.cache, key)
				break
			}
		}
		return nil, errors.New("cache expire")
	}

	for li := lc.l.Front(); li != nil; li = li.Next() {
		if li.Value == key {
			lc.l.MoveToFront(li)
			break
		}
	}
	return itm.Value, nil
}

func (lc *LruCache) String() string {
	var s string
	var i int
	for li := lc.l.Front(); li != nil; li = li.Next() {
		s += fmt.Sprintf("[%d:%s]", i, li.Value)
		i += 1
	}
	return s
}
