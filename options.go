package skiplist

import (
	"time"

	"golang.org/x/exp/constraints"
)

// set options
type setOptions struct {
	setNX   bool
	dueTime time.Time
}
type SetOption interface {
	Apply(o *setOptions)
}

type setOption struct {
	f func(opts *setOptions)
}

func (o *setOption) Apply(opts *setOptions) {
	o.f(opts)
}

func OnNotExist() SetOption {
	return &setOption{f: func(opts *setOptions) {
		opts.setNX = true
	}}
}

func useSetOptions(opts []SetOption) *setOptions {
	v := &setOptions{}
	for _, o := range opts {
		o.Apply(v)
	}
	return v
}

// get options
type getOptions struct {
	getOrDefault bool
}

type GetOption interface {
	Apply(opts *getOptions)
}

type getOption struct {
	f func(opts *getOptions)
}

func (o *getOption) Apply(opts *getOptions) {
	o.f(opts)
}

func GetOrDefault() GetOption {
	return &getOption{f: func(opts *getOptions) {
		opts.getOrDefault = true
	}}
}

func useGetOptions(opts []GetOption) *getOptions {
	v := &getOptions{}
	for _, o := range opts {
		o.Apply(v)
	}
	return v
}

// delete options
type deleteOptions struct {
}
type DeleteOption interface {
	Apply(opts *deleteOptions)
}

// range Options
type rangeOptions[T constraints.Ordered] struct {
	from, to                                T
	includeLowBoundary, includeHighBoundary bool
	hasFrom, hasTo                          bool
}
type RangeOption[T constraints.Ordered] interface {
	Apply(o *rangeOptions[T])
}

type rangeOption[T constraints.Ordered] struct {
	f func(opts *rangeOptions[T])
}

func (o *rangeOption[T]) Apply(opts *rangeOptions[T]) {
	o.f(opts)
}

func From[T constraints.Ordered](from T, includeBoundary bool) RangeOption[T] {
	return &rangeOption[T]{f: func(opts *rangeOptions[T]) {
		opts.from, opts.includeLowBoundary, opts.hasFrom = from, includeBoundary, true
	}}
}

func To[T constraints.Ordered](to T, includeBoundary bool) RangeOption[T] {
	return &rangeOption[T]{f: func(opts *rangeOptions[T]) {
		opts.to, opts.includeHighBoundary, opts.hasTo = to, includeBoundary, true
	}}
}

func useRangeOptions[T constraints.Ordered](opts []RangeOption[T]) *rangeOptions[T] {
	v := &rangeOptions[T]{}
	for _, opt := range opts {
		opt.Apply(v)
	}
	return v
}

// init options
type initOptions struct {
	maxLevels  int
	concurrent bool // default true
}

type InitOption interface {
	Apply(o *initOptions)
}

type initOption struct {
	f func(opts *initOptions)
}

func (o *initOption) Apply(opts *initOptions) {
	o.f(opts)
}

func useInitOptions(opts []InitOption) *initOptions {
	v := &initOptions{concurrent: true, maxLevels: 48}
	for _, o := range opts {
		o.Apply(v)
	}
	return v
}

func Concurrent(c bool) InitOption {
	return &initOption{f: func(opts *initOptions) {
		opts.concurrent = c
	}}
}

func MaxLevels(l int) InitOption {
	return &initOption{f: func(opts *initOptions) {
		opts.maxLevels = l
	}}
}
