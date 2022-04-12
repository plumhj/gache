package gache

import (
	"sync"
	"time"
)

type addable interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | int | float32 | float64 | string
}

const (
	TTLImmortal = int64(0)
)

type Gache[T addable] struct {
	cleanupInterval time.Duration
	items           map[string]item[T]
	lock            sync.RWMutex
	ttlUnit         time.Duration
	OnEviction      func(key string, value T)
}

type item[T addable] struct {
	value T
	ttl   int64
}

func New[T addable](opts ...Option[T]) *Gache[T] {
	ret := &Gache[T]{
		items:   map[string]item[T]{},
		lock:    sync.RWMutex{},
		ttlUnit: time.Second,
	}
	for _, opt := range opts {
		opt(ret)
	}
	if ret.cleanupInterval != 0 {
		go func() {
			for {
				time.Sleep(ret.cleanupInterval)
				go func() {
					targets := make(map[string]T)
					now := time.Now().UnixNano()
					for k, v := range ret.items {
						if now >= v.ttl {
							targets[k] = v.value
						}
					}

					if ret.OnEviction != nil {
						for k, v := range targets {
							ret.OnEviction(k, v)
						}
					}

					ret.lock.Lock()
					for k := range targets {
						delete(ret.items, k)
					}
					ret.lock.Unlock()
				}()
			}
		}()
	}
	return ret
}

func (this *Gache[T]) Set(key string, value T, ttl ...int64) {
	_ttl := TTLImmortal
	if len(ttl) > 0 {
		_ttl = time.Now().UnixNano() + int64(time.Duration(ttl[0])*this.ttlUnit)
	}
	i := item[T]{
		ttl:   _ttl,
		value: value,
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.items[key] = i
}

func (this *Gache[T]) Get(key string) (value T, exist bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	i, ok := this.items[key]
	if ok {
		value = i.value
		exist = true
		if i.ttl != TTLImmortal && time.Now().UnixNano() >= i.ttl {
			exist = false
			value = *new(T)
		}
	}
	return
}

func (this *Gache[T]) Del(key string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.items, key)
}

func (this *Gache[T]) Inc(key string, delta T) T {
	this.lock.Lock()
	defer this.lock.Unlock()

	i, ok := this.items[key]
	if !ok {
		i = item[T]{}
	}
	i.value += delta
	this.items[key] = i
	return i.value
}
