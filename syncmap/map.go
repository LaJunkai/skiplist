package syncmap

type HashAble interface {
	Hash() uint64
}

type SyncMap[KT HashAble, VT any] interface {
	Put(key KT, value VT) (err error)
	Get(key KT, value VT) (err error)
	ContainsKey(key KT) (contain bool, err error)
	Delete(key KT) (err error)
	PutAll(anoMap SyncMap[KT, VT]) (err error)
}
