# Collections
Data structures implements based on Go 1.18+ Generics.  
Data structures included by this project are listed as following.
 - [x] Skiplist
 - [ ] BloomFilter
 - [ ] SyncMap
 - [ ] ConcurrentHashMap

 - [Collections](#collections)
     * [🚀 Install](#---install)
     * [💡 Usage](#---usage)
         + [Skip List](#skip-list)
             - [Basic Usage](#basic-usage)
             - [Instantiation](#instantiation)
             - [Set Options](#set-options)


## 🚀 Install
`go get github.com/LaJunkai/collections`

## 💡 Usage
### Skip List
#### Basic Usage
```go
package main

import (
	"context"
	"fmt"

	"github.com/LaJunkai/collections/skiplist"
)

func main() {
	ctx := context.Background()
	list := skiplist.New[string, string]()

	// set
	if err := list.Set(ctx, "user-0001.name", "John Wick"); err != nil {
		panic(err)
	}

	// get
	value, err := list.Get(ctx, "user-0001.name")
	if err != nil {
		panic(err)
	}
	fmt.Printf("got value: %v\n", value)

	// delete
	deleted, err := list.Delete(ctx, "user-0001.name")
	if err != nil {
		panic(err)
	}
	fmt.Printf("successfully deleted: %v", deleted)

	// range
	err = list.Range(ctx, func(key, value string) bool {
		fmt.Printf("key=%v; value=%v", key, value)
		return true
	})
}

```
#### Instantiation
```go
list := skiplist.New[string, string](
		skiplist.Concurrent(false),
		skiplist.MaxLevels(32),
	)
```  

1. skiplist use mutex to sync operations from different goroutines by default. 
Use `skiplist.Concurrent(false)` to disable the concurrent control.
2. skiplist support levels in range of `1 - 64` (default max level is 48).
Use `skiplist.MaxLevels(n)` to custom the max level limit.

#### Set
```go
err := list.Set(
		ctx, "key-1", "value-1",
		skiplist.OnNotExist(),
	)
```
1. Set method support `OnNotExist()` option. 
With this option passed, an attempt to set an existed key may receive an `ErrDuplicatedKey` error.