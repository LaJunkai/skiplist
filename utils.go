package skiplist

import (
	"crypto/rand"
	"math/big"
	"sync"

	"golang.org/x/exp/constraints"
)

var (
	randomMutex sync.Mutex
)

func randomLevel(maxLevel int, countElement uint64) int {
	randomMutex.Lock()
	defer randomMutex.Unlock()
	
	// Prevent overflow: if countElement is too large, use a safe maximum
	denominator := countElement << 1
	if denominator < countElement || denominator == 0 { // Overflow occurred or zero
		// Use a large but safe value for very large counts
		denominator = 1000000
	}
	denominator += 1
	
	// Ensure denominator is positive for crypto/rand
	if denominator > (1<<63 - 1) {
		denominator = 1000000
	}
	
	// Use crypto/rand for thread-safe random generation
	bigDenominator := big.NewInt(int64(denominator))
	bigV, err := rand.Int(rand.Reader, bigDenominator)
	if err != nil {
		// Fallback to a simple deterministic approach if crypto/rand fails
		return (int(countElement % uint64(maxLevel))) + 1
	}
	
	v := bigV.Uint64()
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
