package fp

func Chunk[T any](arr []T, size int) [][]T {
	var results [][]T

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		results = append(results, arr[i:end])
	}

	return results
}
