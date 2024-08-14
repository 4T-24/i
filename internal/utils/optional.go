package utils

func Optional[K any](k K) *K {
	return &k
}
