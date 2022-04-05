package skiplist

import (
	"golang.org/x/exp/constraints"
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

func max[T constraints.Ordered](values ...T) (maxValue T) {
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

func min[T constraints.Ordered](values ...T) (minValue T) {
	if len(values) > 0 {
		minValue = values[0]
	}
	for _, v := range values {
		if v < minValue {
			minValue = v
		}
	}
	return
}
