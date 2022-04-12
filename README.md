# gache
golang generic cache library

## install

    go get github.com/plumhj/gache

## usage

```go
package main

import (
	"fmt"
	"github.com/plumhj/gache"
	"time"
)

func main() {
	cache := gache.New[string]()
	//no expiration
	cache.Set("KeyA", "ValueA")
	//expire in 10 seconds
	cache.Set("KeyB", "ValueB", 3)
	var str string
	var exist bool
	str, exist = cache.Get("KeyA")
	if exist {
		fmt.Println(str)
	}
	str, exist = cache.Get("KeyB")
	if exist {
		fmt.Println(str)
	}
	time.Sleep(time.Second * 3)
	str, exist = cache.Get("KeyB")
	if !exist {
		fmt.Println("KeyB expired")
	}

	cache2 := gache.New[int]()
	cache2.Inc("KeyA", 1)
	cache2.Inc("KeyA", 1)
	cache2.Inc("KeyA", 1)
	n, exist := cache2.Get("KeyA")
	if exist {
		fmt.Println("KeyA is", n)
	}

	cache3 := gache.New[string](gache.OptionCleanup[string](time.Second))
	cache3.OnEviction = func(key string, value string) {
		fmt.Println("evicted", key, value)
	}
	//KeyA will be evicted in 3 seconds
	cache3.Set("KeyA", "ValueA", 3)
	time.Sleep(time.Second * 5)
}
```