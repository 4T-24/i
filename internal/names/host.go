package names

import (
	"fmt"
	"instancer/internal/env"
)

func GetHost(podName, challengeId, randomId string) string {
	c := env.Get()
	return fmt.Sprintf("%s-%s-%s.%s", podName, challengeId, randomId, c.Domain)
}
