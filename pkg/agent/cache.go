package agent

import "sync"

//go:generate mockery --name=RegisteredResource --output=automock --outpkg=automock
type RegisteredResource[T any] interface {
	Set(T)
	Get() T
}

type Cache[T any] interface {
	Register(name string) RegisteredResource[T]
	Resources() map[string]RegisteredResource[T]
	Clean()
}

type cache[T any] struct {
	resources map[string]RegisteredResource[T]
}

func NewCache[T any]() Cache[T] {
	return &cache[T]{
		resources: map[string]RegisteredResource[T]{},
	}
}

func (c *cache[T]) Register(name string) RegisteredResource[T] {
	r := &resource[T]{}
	c.resources[name] = r
	return r
}

func (c *cache[T]) Resources() map[string]RegisteredResource[T] {
	return c.resources
}

func (c *cache[T]) Clean() {
	c.resources = map[string]RegisteredResource[T]{}
}

type resource[T any] struct {
	m   sync.Mutex
	val T
}

func (r *resource[T]) Set(val T) {
	r.m.Lock()
	defer r.m.Unlock()

	r.val = val
}

func (r *resource[T]) Get() T {
	r.m.Lock()
	defer r.m.Unlock()

	return r.val
}
