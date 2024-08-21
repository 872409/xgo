package fp

import "errors"

var ERR_NOT_Found = errors.New("NotFound")

func FilterFirstMust[T any](items []T, fn func(index int, item T) bool) T {
	result, _ := FilterFirst[T](items, fn)
	return result
}

func FilterFirst[T any](items []T, fn func(index int, item T) bool) (T, error) {
	var zeroValue T

	if items == nil || len(items) == 0 {
		return zeroValue, ERR_NOT_Found
	}

	for index, value := range items {
		if fn(index, value) {
			return value, nil
		}
	}

	return zeroValue, ERR_NOT_Found
}

func filter[T any](items []T, condition func(index int, item T) bool, filterMatch bool) (result []T) {
	if items == nil || len(items) == 0 {
		return result
	}

	for index, value := range items {
		if condition(index, value) == filterMatch {
			result = append(result, value)
		}
	}

	return result
}

func Filter[T any](items []T, condition func(index int, item T) bool) []T {
	return filter(items, condition, true)
}

func FilterNot[T any](items []T, condition func(index int, item T) bool) []T {
	return filter(items, condition, false)
}

func RemoveDuplicate[T comparable](slc []T) []T {
	var result []T
	tempMap := map[T]byte{}

	for _, e := range slc {
		if _, found := tempMap[e]; !found {
			tempMap[e] = 1
			result = append(result, e)
		}
	}
	return result
}

//
//func Map[T any](items []T, fn func(index int, item T) bool) error {
//	for index, value := range items {
//		if fn(index, value) {
//			return nil
//		}
//	}
//	return ERR_NOT_Found
//}

func Map[T any, P any](items []T, mapFunc func(index int, item T) P) []P {
	var p []P
	for idx, value := range items {
		p = append(p, mapFunc(idx, value))
	}
	return p
}

func Includes[T comparable](items []T, item T) bool {
	for _, value := range items {
		if value == item {
			return true
		}
	}
	return false
}
func IncludesCondition[T comparable](items []T, condition func(index int, item T) bool) bool {
	for idx, value := range items {
		if condition(idx, value) {
			return true
		}
	}
	return false
}

func OfOne[T comparable](items []T, item ...T) bool {
	for _, value := range items {
		if Includes(item, value) {
			return true
		}
	}
	return false
}
