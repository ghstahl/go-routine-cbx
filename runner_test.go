package cbx

import (
	"testing"
)

type SomeContainer struct {
	Name string
}

func TestNewHandlerContainer_Cache(t *testing.T) {
	container := NewHandlerContainer(nil)

	cache := container.GetCache()
	if cache != nil {
		t.Fail()
	}
	v := SomeContainer{
		Name: "tabby",
	}
	container.UpsertCache(v)
	cache = container.GetCache()
	if cache == nil {
		t.Fail()
	}
	var someC SomeContainer = cache.(SomeContainer)
	if someC.Name != v.Name {
		t.Fail()
	}

	t.Log(cache)
}
func TestNewHandlerContainer_CachePointer(t *testing.T) {
	container := NewHandlerContainer(nil)

	cache := container.GetCache()
	if cache != nil {
		t.Fail()
	}
	v := SomeContainer{
		Name: "tabby",
	}
	container.UpsertCache(&v)
	cache = container.GetCache()
	if cache == nil {
		t.Fail()
	}
	var someC *SomeContainer = cache.(*SomeContainer)
	if someC.Name != v.Name {
		t.Fail()
	}

	t.Log(cache)
}
