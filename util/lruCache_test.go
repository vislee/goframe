package util

import (
	"strings"
	"testing"
	"time"
)

// liwq liwq4 liwq3 liwq2
func TestUpdate(t *testing.T) {
	cache := NewLruCache(3, 5)
	cache.Init()
	cache.Update("liwq", 123)
	cache.Update("liwq2", 456)
	cache.Update("liwq3", 789)
	cache.Update("liwq4", 888)
	// liwq3 liwq liwq4
	it, err := cache.Search("liwq3")
	if err != nil || it != 789 {
		t.Error("update search error")
	}
	s := cache.String()
	if !strings.HasPrefix(s, "[0:liwq3]") {
		t.Error("update error")
		t.Error(s)
	}
	cache.Update("liwq5", 999)
	//liwq3 liwq5 liwq
	s = cache.String()
	if !strings.HasPrefix(s, "[0:liwq3]") {
		t.Error("update error")
		t.Error(s)
	}
}

func TestSearchOk(t *testing.T) {
	cache := NewLruCache(3, 5)
	cache.Init()
	cache.Update("liwq", 123)
	cache.Update("liwq2", 456)
	cache.Update("liwq3", 789)
	cache.Update("liwq4", 888)
	_, err := cache.Search("liwq3")
	if err != nil {
		t.Error(err.Error())
		t.Error(cache.String())
	}
}

func TestSearchMiss(t *testing.T) {
	cache := NewLruCache(3, 2)
	cache.Init()
	cache.Update("liwq", 123)
	cache.Update("liwq2", 456)
	cache.Update("liwq3", 789)
	cache.Update("liwq4", 888)
	_, err := cache.Search("liwq2")
	if err != nil {
		t.Error(err.Error())
		t.Error(cache.String())
	}
}

func TestSearchExpire(t *testing.T) {
	cache := NewLruCache(3, 2)
	cache.Init()
	cache.Update("liwq", 123)
	time.Sleep(3 * time.Second)
	_, err := cache.Search("liwq")
	if err == nil {
		t.Error("lru expire error")
		t.Error(cache.String())
	}
}

type user struct {
	name string
	age  int
}

func TestUpdateStruct(t *testing.T) {
	cache := NewLruCache(3, 2)
	cache.Init()
	cache.Update("liwq", user{"liwq", 23})
	cache.Update("liwq1", user{"liwq1", 34})
	cache.Update("liwq2", user{"liwq2", 45})
	it, err := cache.Search("liwq2")
	if err != nil {
		t.Error("update struct error")
	}
	u, e := it.(user)
	if !e {
		t.Error("update struct interface error")
	}
	if u.name != "liwq2" || u.age != 45 {
		t.Error("update struct value error")
	}
}
