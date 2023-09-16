package cache

import (
	"sync"
	"time"

	"github.com/pPrecel/cloudagent/pkg/types"
)

type GardenerCache Cache[*types.ShootList]
type GardenerRegisteredResource RegisteredResource[*types.ShootList]

func NewGardenerCache() GardenerCache {
	return NewCache[*types.ShootList]()
}

type ServerCache struct {
	GardenerCache GardenerCache
	GeneralError  error
}

func (sc *ServerCache) GetGardenerCache() GardenerCache {
	return sc.GardenerCache
}

func (sc *ServerCache) GetGeneralError() error {
	return sc.GeneralError
}

//go:generate mockery --name=RegisteredResource --output=automock --outpkg=automock
type RegisteredResource[T any] interface {
	Set(T, error)
	Get() Value[T]
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

type Value[T any] struct {
	Error error
	Time  time.Time
	Value T
}

type resource[T any] struct {
	m   sync.Mutex
	val Value[T]
}

// Set - set value with time or error
func (r *resource[T]) Set(val T, err error) {
	r.m.Lock()
	defer r.m.Unlock()

	// set error only to not override existing last set value
	if err != nil {
		r.val.Error = err
		return
	}

	// create new object to override last error
	r.val = Value[T]{
		Value: val,
		Error: nil,
		Time:  time.Now(),
	}
}

// Get - returns latest Value object
func (r *resource[T]) Get() Value[T] {
	r.m.Lock()
	defer r.m.Unlock()

	return r.val
}
