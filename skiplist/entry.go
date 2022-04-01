package skiplist

import "golang.org/x/exp/constraints"

type ReadOnlyEntry[KT constraints.Ordered, VT any] interface {
	Key() KT
	Value() VT
}

type readonlyEntry[KT constraints.Ordered, VT any] struct {
	key   KT
	value VT
}

func (r *readonlyEntry[KT, VT]) Key() KT {
	return r.key
}

func (r *readonlyEntry[KT, VT]) Value() VT {
	return r.value
}

type entry[KT constraints.Ordered, VT any] struct {
	readonlyEntry[KT, VT]
	levels []*entry[KT, VT]
	prev   *entry[KT, VT]
}

func newEntry[KT constraints.Ordered, VT any](key KT, value VT, levels int) *entry[KT, VT] {
	return &entry[KT, VT]{
		readonlyEntry: readonlyEntry[KT, VT]{key: key, value: value},
		levels:        make([]*entry[KT, VT], levels),
	}
}

func newEmptyEntry[KT constraints.Ordered, VT any](levels int) *entry[KT, VT] {
	return &entry[KT, VT]{
		levels: make([]*entry[KT, VT], levels),
	}
}
