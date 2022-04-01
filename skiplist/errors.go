package skiplist

import "errors"

var (
	ErrKeyNotFound   = errors.New("the specified key is not exist")
	ErrDuplicatedKey = errors.New("the specified key has already existed")
)
