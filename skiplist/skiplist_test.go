package skiplist

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

var (
	key   []int
	value []int
)

func init() {
	key, value = make([]int, 5000000), make([]int, 5000000)
	for i := 0; i < 5000000; i++ {
		key[i] = rand.Intn(5000000)
		value[i] = i
	}
}

func TestSet(t *testing.T) {
	var ctx = context.Background()
	list := New[string, string]()
	if err := list.Set(ctx, "key-99", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-98", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-97", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-999", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-99", "key-99", SetNX()); err != ErrDuplicatedKey {
		t.Error("expect ErrDuplicatedKey")
	}
}

func TestSkipList_Get(t *testing.T) {
	var ctx = context.Background()
	list := New[int, int]()
	for i := 0; i < 1000000; i++ {
		if err := list.Set(ctx, key[i], value[i]); err != nil {
			t.Error(err)
		}
		if v, err := list.Get(ctx, key[i]); v != value[i] || err != nil {
			t.Error("list.Get return an err or a wrong value")
		}
	}
}

func TestSkipList_Delete(t *testing.T) {
	var ctx = context.Background()
	list := New[string, string]()
	if err := list.Set(ctx, "key-99", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-98", "key-99"); err != nil {
		t.Error(err)
	}
	if err := list.Set(ctx, "key-97", "abc"); err != nil {
		t.Error(err)
	}
	if deleted, err := list.Delete(ctx, "key-97"); err != nil {
		t.Error(err)
	} else if !deleted {
		t.Errorf("should delete key %v", "key-97")
	}
	if v, err := list.Get(ctx, "key-97"); err != ErrKeyNotFound {
		fmt.Println(v)
		t.Error("except ErrKeyNotFound")
	}
}

func BenchmarkSkipList_Set(b *testing.B) {
	list := New[int, int](Concurrent(true))
	ctx := context.Background()
	for i := b.N; i > 0; i-- {
		if err := list.Set(ctx, key[i], value[i]); err != nil {
			b.Error(err)
		}
	}
	fmt.Println(b.N)
}
