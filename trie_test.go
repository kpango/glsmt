package glsmt

import (
	"testing"
)

func TestTrie(t *testing.T) {
	tr := NewTrie[string](0)

	keys := []string{"apple", "app", "bat", "ball"}
	for _, key := range keys {
		k := key
		tr.Insert(key, &k)
	}

	// tr.Range(func(key string, value *string) bool {
	// t.Logf("key: %s, value: %s", key, *value)
	// return true
	// })

	for _, key := range keys {
		val, ok := tr.Get(key)
		if !ok || *val != key {
			t.Errorf("Get failed for key %s, got %v", key, val)
		}
	}

	if val, ok := tr.Get("nokey"); ok {
		t.Errorf("Get should fail for non-existing key, got %v", *val)
	}

	tr.Delete("app")
	if val, ok := tr.Get("app"); ok {
		t.Errorf("Get should fail for deleted key, got %v", *val)
	}

	tr.Delete("nokey")

	val, ok := tr.Get("apple")
	if !ok || *val != "apple" {
		t.Errorf("Get failed for key apple after deleting app, got %v", *val)
	}
}
