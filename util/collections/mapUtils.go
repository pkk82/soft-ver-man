package collections

func GroupBy[T any, K comparable](values []T, classifier func(T) K) map[K][]T {
	grouped := make(map[K][]T)
	for _, v := range values {
		key := classifier(v)
		grouped[key] = append(grouped[key], v)
	}
	return grouped
}

func GroupByAndCollect[T, S any, K comparable](values []T, classifier func(T) K, collector func([]T) S) map[K]S {
	groupedAndCollected := make(map[K]S)
	grouped := GroupBy(values, classifier)
	for k, v := range grouped {
		groupedAndCollected[k] = collector(v)
	}
	return groupedAndCollected
}
