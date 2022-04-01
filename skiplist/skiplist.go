package skiplist

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/exp/constraints"
)

const (
	maxSupportedLevel = 64
)

type SkipList[KT constraints.Ordered, VT any] interface {
	Get(ctx context.Context, key KT, opts ...GetOption) (value VT, err error)
	Set(ctx context.Context, key KT, value VT, opts ...SetOption) (err error)
	Range(ctx context.Context, f func(key KT, value VT) bool, opts ...RangeOption[KT]) (err error)
	Delete(ctx context.Context, key KT, opts ...DeleteOption) (deleted bool, err error)
	Size() uint64
	CurrentMaxLevel() int
	debugLevels()
}

type skipList[KT constraints.Ordered, VT any] struct {
	head            *entry[KT, VT]
	rwMutex         sync.RWMutex
	initOpts        *initOptions
	countElement    uint64
	maxLevel        int
	currentMaxLevel int
}

func (list *skipList[KT, VT]) Set(ctx context.Context, key KT, value VT, opts ...SetOption) (err error) {
	if list.initOpts.concurrent {
		list.rwMutex.Lock()
		defer list.rwMutex.Unlock()
	}
	return list.set(ctx, key, value, opts...)
}

func (list *skipList[KT, VT]) Get(ctx context.Context, key KT, opts ...GetOption) (value VT, err error) {
	if list.initOpts.concurrent {
		list.rwMutex.RLock()
		defer list.rwMutex.RUnlock()
	}
	return list.get(ctx, key, opts...)
}

func (list *skipList[KT, VT]) Delete(ctx context.Context, key KT, opts ...DeleteOption) (deleted bool, err error) {
	if list.initOpts.concurrent {
		list.rwMutex.Lock()
		defer list.rwMutex.Unlock()
	}
	return list.delete(ctx, key, opts...)
}
func (list *skipList[KT, VT]) debugLevels() {
	current := list.head
	for current != nil {
		fmt.Println(strings.Repeat("*", len(current.levels)))
		current = current.levels[0]
	}
}

func (list *skipList[KT, VT]) Range(ctx context.Context,
	f func(key KT, value VT) bool, opts ...RangeOption[KT]) (err error) {
	if list.initOpts.concurrent {
		list.rwMutex.RLock()
		defer list.rwMutex.RUnlock()
	}
	current := list.head.levels[0]
	o := useRangeOptions(opts)
	if o.hasFrom {
		prevList, _, err := list.find(ctx, o.from, list.initOpts.maxLevels, false)
		if err != nil {
			return err
		}
		current = prevList[0].levels[0]

		if current == nil {
			return nil
		}
		if current.key < o.from || (current.key == o.from && !o.includeLowBoundary) {
			current = current.levels[0]
		}
	}

	for current != nil {
		if o.hasTo {
			if current.key > o.to {
				return nil
			}
			if current.key == o.to && !o.includeHighBoundary {
				return nil
			}
		}
		if !f(current.key, current.value) || ctx.Err() != nil {
			break
		}
		current = current.levels[0]
	}
	return nil
}

func (list *skipList[KT, VT]) get(ctx context.Context, key KT, opts ...GetOption) (value VT, err error) {
	prevEntries, findAt, err := list.find(ctx, key, list.initOpts.maxLevels, true)
	if err != nil {
		return
	}
	target := prevEntries[findAt].levels[findAt]
	if target == nil || target.key != key {
		return value, ErrKeyNotFound
	}
	return target.value, err
}

func (list *skipList[KT, VT]) set(ctx context.Context, key KT, value VT, opts ...SetOption) (err error) {
	o := useSetOptions(opts)

	rLevel := randomLevel(list.initOpts.maxLevels, list.countElement)
	list.currentMaxLevel = max(list.currentMaxLevel, rLevel)
	prevEntries, _, err := list.find(ctx, key, rLevel, false)
	if err != nil {
		return err
	}

	if e := prevEntries[0].levels[0]; e != nil && e.key == key {
		if o.setNX {
			return ErrDuplicatedKey
		}
		e.value = value
		return nil
	}

	// insert the new entry
	e := newEntry(key, value, rLevel+1)
	for i := 0; i <= rLevel; i++ {
		prev := prevEntries[i]
		next := prev.levels[i]
		prev.levels[i], e.levels[i] = e, next
	}
	list.countElement++
	return nil
}

func (list *skipList[KT, VT]) delete(ctx context.Context, key KT, opts ...DeleteOption) (deleted bool, err error) {
	prevEntries, _, err := list.find(ctx, key, list.initOpts.maxLevels, false)

	if e := prevEntries[0].levels[0]; e != nil && e.key == key {
		for level, next := range e.levels {
			prevEntries[level].levels[level] = next
		}
		list.countElement--
		return true, nil
	}
	return false, err

}

func (list *skipList[KT, VT]) find(ctx context.Context,
	key KT, level int, returnOnFind bool) (prevEntries []*entry[KT, VT], findAt int, err error) {
	prevEntries = make([]*entry[KT, VT], level+1)
	prev := list.head
	for i := list.currentMaxLevel; i >= 0; i-- {
		current := prev.levels[i]
		for {
			if current == nil || current.key >= key {
				if i <= level {
					prevEntries[i] = prev
					if returnOnFind && current != nil && current.key == key {
						return prevEntries, i, nil
					}
				}
				break
			}
			prev, current = current, current.levels[i]
		}
		if ctx.Err() != nil {
			return prevEntries, 0, ctx.Err()
		}
	}
	return prevEntries, 0, nil
}

func (list *skipList[KT, VT]) Size() uint64 {
	return list.countElement
}

func (list *skipList[KT, VT]) CurrentMaxLevel() int {
	return list.currentMaxLevel
}

func New[KT constraints.Ordered, VT any](opts ...InitOption) SkipList[KT, VT] {
	o := useInitOptions(opts)
	if o.maxLevels <= 0 || o.maxLevels > maxSupportedLevel {
		o.maxLevels = 48
	}
	return &skipList[KT, VT]{
		head:     newEmptyEntry[KT, VT](o.maxLevels + 1),
		initOpts: o,
	}
}
