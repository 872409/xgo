package fp

import (
	"math/rand"
	"time"
)

func GroupBy[T any, K comparable](items []T, getProperty func(T) K) map[K][]T {
	grouped := make(map[K][]T)

	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}

func SliceShuffle(a []any) {
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func SliceShuffleRand(r *rand.Rand, a []any) {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().Unix()))
	}
	r.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}
