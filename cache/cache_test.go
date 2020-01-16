package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewCache(1)
	defer c.Stop()

	c.Set("foo", 123, 0)
	r, ok := c.Get("foo")
	if !ok || r.(int) != 123 {
		t.Errorf("err on get: got (%d, %t) <> want (%d, %t)", r, ok, 123, true)
	}

	c.Set("foo", 321, time.Second)
	r, ok = c.Get("foo")
	if !ok || r.(int) != 321 {
		t.Errorf("err on get: got (%d, %t) <> want (%d, %t)", r, ok, 321, true)
	}

	c.Remove("foo")

	_, ok = c.Get("foo")
	if ok {
		t.Error("err on remove")
	}
}

func TestCacheKeys(t *testing.T) {
	c := NewCache(1)

	c.Set("foo", 123, 0)
	c.Set("bar", 321, 0)
	c.Set("foo", 666, 0)

	keys := c.Keys()

	if len(keys) != 2 {
		t.Error("err on getting keys - wrong number of keys")
	}

	if !(keys[0] == "foo" && keys[1] == "bar") && !(keys[1] == "foo" && keys[0] == "bar") {
		t.Errorf("err on getting keys: got {%s, %s} <> want {%s, %s}", keys[0], keys[1], "foo", "bar")
	}
}

func TestCacheRemovingExpired(t *testing.T) {
	c := NewCache(2)

	c.Set("foo", 321, time.Second)
	time.Sleep(time.Second * 4)

	_, ok := c.Get("foo")
	if ok {
		t.Error("err on removing expired key")
	}
}
