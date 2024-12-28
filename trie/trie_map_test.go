package trie

import (
	"fmt"
	"testing"
)

// set get 预期是设置进去， 也能读出来
func Test_TrimeMap_SetGet(t *testing.T) {

	tm := New[string]()
	max := 1000

	for i := 1; i < max; i++ {

		key := fmt.Sprint(i)

		tm.Set(key, key)
		val := tm.Get(key)
		if key != val {
			t.Errorf("expected %s, got %s", key, val)
		}
	}
}

// HasPrefix 找到
func Test_TrieMap_HasPrefix(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		if !tm.HasPrefix(key[:i]) {
			t.Errorf("expected true for prefix %s", key[:i])
		}
	}
}

// HasPrefix 找不到
func Test_TrieMap_HasPrefix_notFound(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		if !tm.HasPrefix(key[:i]) {
			t.Errorf("expected true for prefix %s", key[:i])
		}
	}

	if tm.HasPrefix("/ha") {
		t.Errorf("expected false for prefix /ha")
	}
}

func Test_TrieMap_TryGet_notFound(t *testing.T) {
	tm := New[string]()
	key := "/hello/world"
	tm.Set("/hello", "1")
	tm.Set("/hello/world", "1")
	for i := 1; i < len(key); i++ {

		if !tm.HasPrefix(key[:i]) {
			t.Errorf("expected true for prefix %s", key[:i])
		}
	}
	_, ok := tm.TryGet("/ha")
	if ok {
		t.Errorf("expected false for /ha")
	}
	_, ok = tm.TryGet("/he")
	if ok {
		t.Errorf("expected false for /he")
	}
}

func Test_TrieMap_Delete(t *testing.T) {

	tm := New[string]()
	max := 1000

	for i := 1; i < max; i++ {

		key := fmt.Sprint(i)

		tm.Set(key, key)
		val := tm.Get(key)
		if key != val {
			t.Errorf("expected %s, got %s", key, val)
		}
		tm.Delete(key)
		val, ok := tm.TryGet(key)
		if ok {
			t.Errorf("expected false for key %s", key)
		}
		if val != "" {
			t.Errorf("expected empty string, got %s", val)
		}
	}

	key := fmt.Sprint(max + 1)
	tm.Delete(key)
	val, ok := tm.TryGet(key)
	if ok {
		t.Errorf("expected false for key %s", key)
	}
	if val != "" {
		t.Errorf("expected empty string, got %s", val)
	}
}

// 删除长的
func Test_TrieMap_Delete2(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Delete("/12")
	if tm.Get("/12") != "" {
		t.Errorf("expected empty string for /12, got %s", tm.Get("/12"))
	}
	if tm.Get("/1") != "/1" {
		t.Errorf("expected /1, got %s", tm.Get("/1"))
	}
}

// 删除短的
func Test_TrieMap_Delete3(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Delete("/1")
	if tm.Get("/12") != "/12" {
		t.Errorf("expected /12, got %s", tm.Get("/12"))
	}
	if tm.Get("/1") != "" {
		t.Errorf("expected empty string for /1, got %s", tm.Get("/1"))
	}
}

// 删除带中文
func Test_TrieMap_Delete4(t *testing.T) {

	tm := New[string]()

	tm.Set("中", "中")
	tm.Set("中国", "中国")
	tm.Delete("中")
	if tm.Get("中国") != "中国" {
		t.Errorf("expected 中国, got %s", tm.Get("中国"))
	}
	if tm.Get("中") != "" {
		t.Errorf("expected empty string for 中, got %s", tm.Get("中"))
	}
}

func Test_TrieMap_Delete5(t *testing.T) {

	tm := New[string]()

	tm.Set("中", "中")
	tm.Set("中国", "中国")
	tm.Delete("中")
	tm.Delete("中国")
	if tm.Get("中国") != "" {
		t.Errorf("expected empty string for 中国, got %s", tm.Get("中国"))
	}
	if tm.Get("中") != "" {
		t.Errorf("expected empty string for 中, got %s", tm.Get("中"))
	}
}

func Test_TrieMap_Delete6(t *testing.T) {

	tm := New[string]()

	tm.Set("/1", "/1")
	tm.Set("/12", "/12")
	tm.Set("/13", "/13")
	tm.Delete("/12")
	if tm.Get("/12") != "" {
		t.Errorf("expected empty string for /12, got %s", tm.Get("/12"))
	}
	if tm.Get("/1") != "/1" {
		t.Errorf("expected /1, got %s", tm.Get("/1"))
	}
	if tm.Get("/13") != "/13" {
		t.Errorf("expected /13, got %s", tm.Get("/13"))
	}
}
