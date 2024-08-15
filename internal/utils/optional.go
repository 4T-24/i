package utils

import "fmt"

func Optional[K any](k K) *K {
	return &k
}

func SprintPtr[K any](k *K) string {
	if k == nil {
		return ""
	}
	return fmt.Sprint(*k)
}
