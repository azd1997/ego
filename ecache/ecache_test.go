package ecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

// 模拟数据库(源)
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestECache_Get(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	ecache := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if item, err := ecache.Get(k); err != nil || item.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, err := ecache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if item, err := ecache.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", item)
	}
}


