package skiplist

import (
	"math/rand"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func init() {

}

func randomLevel(maxLevel int, countElement uint64) int {
	v := random.Uint64() % (countElement + 1)
	for i := 1; i <= maxLevel; i++ {
		if v&1 == 0 {
			return i
		}
		v >>= 1
	}
	return maxLevel
}

func max[T int](values ...T) (maxValue T) {
	if len(values) > 0 {
		maxValue = values[0]
	}
	for _, v := range values {
		if v > maxValue {
			maxValue = v
		}
	}
	return
}
