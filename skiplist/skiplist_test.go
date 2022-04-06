package skiplist

import (
	"math/rand"
	"strconv"
	"testing"
)

var (
	keys           []string
	values         []int
	scale          = 1000000
	maxKey, minKey = "", "{"
	readonlyList   SkipList[string, int]
)

func init() {
	readonlyList = New[string, int]()
	keys, values = make([]string, scale), make([]int, scale)
	for i := 0; i < scale; i++ {
		keys[i] = strconv.Itoa(rand.Intn(scale))
		values[i] = i
		maxKey = max(keys[i], maxKey)
		minKey = min(keys[i], minKey)
		_ = readonlyList.Set(keys[i], values[i])
	}
}

func TestNew(t *testing.T) {
	list := New[float64, string]()
	if list == nil {
		t.Error("expected the list instance, got nil")
	}
}

func TestSkipList_Set(t *testing.T) {
	list := New[string, int](MaxLevels(18))
	for i := 0; i < scale; i++ {
		if err := list.Set(keys[i], values[i]); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSkipList_Get(t *testing.T) {
	list := New[string, int]()
	for i := 0; i < scale; i++ {
		if err := list.Set(keys[i], values[i]); err != nil {
			t.Fatal(err)
		}
		value, err := list.Get(keys[i])
		if err != nil {
			t.Fatal(err)
		}
		if value != values[i] {
			t.Fatal("Get method return the wrong value")
		}
	}
}

func TestSkipList_Get2(t *testing.T) {
	_, err := readonlyList.Get("{{{{")
	if err != ErrKeyNotFound {
		t.Fatal("expected ErrKeyNotFound")
	}
}

func TestSkipList_Delete(t *testing.T) {
	list := New[string, int]()
	for i := 0; i < scale; i++ {
		if err := list.Set(keys[i], values[i]); err != nil {
			t.Fatal(err)
		}
		if i%4 == 0 {
			deleted := list.Delete(keys[i])
			if !deleted {
				t.Fatalf("expect delete the key [%v]", keys[i])
			}
		}
	}
	deleted := list.Delete("{{{")
	if deleted {
		t.Error("expected deleted=false, bot deleted=true")
	}
}

func TestSkipList_Range(t *testing.T) {
	list := New[string, int]()
	kvs := []struct {
		key   string
		value int
	}{{key: "user-001", value: 1}, {key: "user-022", value: 3}, {key: "user-098", value: 87}}
	for _, item := range kvs {
		_ = list.Set(item.key, item.value)
	}
	index := 0
	list.Range(func(key string, value int) bool {
		if kv := kvs[index]; key != kv.key || value != kv.value {
			t.Fatalf("range disordered, expected %v=%v, got %v=%v", kv.key, kv.value, key, value)
		}
		index += 1
		return true
	})

}

func TestSkipList_Pop(t *testing.T) {
	list := New[string, int]()
	_ = list.Set("user-07", 1)
	_ = list.Set("user-99", 2)
	key, value, err := list.Pop()
	if err != nil {
		t.Error(err)
	}
	if key != "user-99" || value != 2 {
		t.Fatal("wrong key-value popped")
	}
	key, value, err = list.Pop()
	if err != nil {
		t.Error(err)
	}
	if key != "user-07" || value != 1 {
		t.Fatal("wrong key-value popped")
	}
	_, _, err = list.Pop()
	if err != ErrNoEntries {
		t.Error("expected ErrNoEntries")
	}
}

func TestSkipList_LPop(t *testing.T) {
	list := New[string, int]()
	_ = list.Set("user-07", 1)
	_ = list.Set("user-99", 2)
	key, value, err := list.LPop()
	if err != nil {
		t.Error(err)
	}
	if key != "user-07" || value != 1 {
		t.Fatal("wrong key-value popped")
	}
	key, value, err = list.LPop()
	if err != nil {
		t.Error(err)
	}
	if key != "user-99" || value != 2 {
		t.Fatal("wrong key-value popped")
	}
	_, _, err = list.LPop()
	if err != ErrNoEntries {
		t.Error("expected ErrNoEntries")
	}
}

func TestSkipList_Max(t *testing.T) {
	key, _, err := readonlyList.Max()
	if err != nil {
		t.Fatal(err)
	}
	if key != maxKey {
		t.Fatal("Max method return the wrong value")
	}
}

func TestSkipList_Min(t *testing.T) {
	key, _, err := readonlyList.Min()
	if err != nil {
		t.Fatal(err)
	}
	if key != minKey {
		t.Fatal("Min method return the wrong value")
	}
}

func TestSkipList_Size(t *testing.T) {
	list := New[int, int]()
	_ = list.Set(1, 2)
	if list.Size() != 1 {
		t.Error("expected size of 1")
	}
	_ = list.Set(1, 3)
	if list.Size() != 1 {
		t.Error("expected size of 1")
	}
	_ = list.Delete(1)
	if list.Size() != 0 {
		t.Error("expected size of 0")
	}
}
