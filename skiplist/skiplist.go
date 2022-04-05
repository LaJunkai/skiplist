package skiplist

import (
	"sync"

	"golang.org/x/exp/constraints"
)

const (
	maxSupportedLevel = 64
)

type SkipList[KT constraints.Ordered, VT any] interface {
	Get(key KT, opts ...GetOption) (value VT, err error)
	Set(key KT, value VT, opts ...SetOption) (err error)
	Range(f func(key KT, value VT) bool, opts ...RangeOption[KT]) (err error)
	Delete(key KT, opts ...DeleteOption) (deleted bool, err error)
	Pop() (key KT, value VT, err error)
	LPop() (key KT, value VT, err error)
	Max() (key KT, value VT, err error)
	Min() (key KT, value VT, err error)
	Size() uint64
}

type skipList[KT constraints.Ordered, VT any] struct {
	head            *entry[KT, VT]
	tail            *entry[KT, VT]
	rwMutex         sync.RWMutex
	initOpts        *initOptions
	countElement    uint64
	maxLevel        int
	currentMaxLevel int
}

// Set method is used to put kv pairs into the Skiplist it support `OnNotExist()` option.
// With this option enabled, an attempt to commit set operation on an existed key may receive an `ErrDuplicatedKey` error.
// err := list.Set(
//   "key-1", "value-1",
//   skiplist.OnNotExist(),
// )
func (list *skipList[KT, VT]) Set(key KT, value VT, opts ...SetOption) (err error) {
	if list.initOpts.concurrent {
		list.rwMutex.Lock()
		defer list.rwMutex.Unlock()
	}
	return list.set(key, value, opts...)
}

func (list *skipList[KT, VT]) Get(key KT, opts ...GetOption) (value VT, err error) {
	if list.initOpts.concurrent {
		list.rwMutex.RLock()
		defer list.rwMutex.RUnlock()
	}
	return list.get(key, opts...)
}

func (list *skipList[KT, VT]) Delete(key KT, opts ...DeleteOption) (deleted bool, err error) {
	if list.initOpts.concurrent {
		list.rwMutex.Lock()
		defer list.rwMutex.Unlock()
	}
	return list.delete(key, opts...)
}

func (list *skipList[KT, VT]) Range(f func(key KT, value VT) bool, opts ...RangeOption[KT]) (err error) {
	if list.initOpts.concurrent {
		list.rwMutex.RLock()
		defer list.rwMutex.RUnlock()
	}
	current := list.head.levels[0]
	o := useRangeOptions(opts)
	if o.hasFrom {
		prevList, _, err := list.find(o.from, list.initOpts.maxLevels, false)
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
		if !f(current.key, current.value) {
			break
		}
		current = current.levels[0]
	}
	return nil
}

func (list *skipList[KT, VT]) Max() (key KT, value VT, err error) {
	return list.max()
}

func (list *skipList[KT, VT]) Min() (key KT, value VT, err error) {
	return list.min()
}

func (list *skipList[KT, VT]) Pop() (key KT, value VT, err error) {
	return list.pop()
}

func (list *skipList[KT, VT]) LPop() (key KT, value VT, err error) {
	return list.lPop()
}

func (list *skipList[KT, VT]) get(key KT, opts ...GetOption) (value VT, err error) {
	prevEntries, findAt, err := list.find(key, list.initOpts.maxLevels, true)
	if err != nil {
		return
	}
	target := prevEntries[findAt].levels[findAt]
	if target == nil || target.key != key {
		return value, ErrKeyNotFound
	}
	return target.value, err
}

func (list *skipList[KT, VT]) set(key KT, value VT, opts ...SetOption) (err error) {
	o := useSetOptions(opts)

	rLevel := randomLevel(list.initOpts.maxLevels, list.countElement)
	list.currentMaxLevel = max(list.currentMaxLevel, rLevel)
	prevEntries, _, err := list.find(key, rLevel, false)
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
	if next := e.levels[0]; next == nil {
		list.tail = e
	} else {
		next.prev = e
	}
	e.prev = prevEntries[0]
	list.countElement++
	return nil
}

func (list *skipList[KT, VT]) delete(key KT, opts ...DeleteOption) (deleted bool, err error) {
	prevEntries, _, err := list.find(key, list.initOpts.maxLevels, false)

	if e := prevEntries[0].levels[0]; e != nil && e.key == key {
		for level, next := range e.levels {
			prevEntries[level].levels[level] = next
		}
		if prev := prevEntries[0]; prev.levels[0] == nil {
			list.tail = prevEntries[0]
		} else {
			prev.levels[0].prev = prev
		}
		list.countElement--
		return true, nil
	}
	return false, err
}

func (list *skipList[KT, VT]) pop() (key KT, value VT, err error) {
	key, value, err = list.max()
	if err != nil {
		return
	}
	_, err = list.delete(key)
	return
}

func (list *skipList[KT, VT]) lPop() (key KT, value VT, err error) {
	key, value, err = list.min()
	if err != nil {
		return
	}
	_, err = list.delete(key)
	return
}

func (list *skipList[KT, VT]) max() (key KT, value VT, err error) {
	if list.countElement == 0 {
		err = ErrNoEntries
		return
	}
	return list.tail.key, list.tail.value, nil
}

func (list *skipList[KT, VT]) min() (key KT, value VT, err error) {
	if list.countElement == 0 {
		err = ErrNoEntries
		return
	}
	return list.head.levels[0].key, list.head.levels[0].value, nil
}

func (list *skipList[KT, VT]) find(key KT, level int, returnOnFind bool) (prevEntries []*entry[KT, VT], findAt int, err error) {
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
	}
	return prevEntries, 0, nil
}

func (list *skipList[KT, VT]) Size() uint64 {
	return list.countElement
}

func New[KT constraints.Ordered, VT any](opts ...InitOption) SkipList[KT, VT] {
	o := useInitOptions(opts)
	if o.maxLevels <= 0 || o.maxLevels > maxSupportedLevel {
		o.maxLevels = 48
	}
	head := newEmptyEntry[KT, VT](o.maxLevels + 1)
	return &skipList[KT, VT]{
		head:     head,
		tail:     head,
		initOpts: o,
	}
}
