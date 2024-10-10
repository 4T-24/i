package names

import (
	"fmt"
	"instancer/internal/env"
)

func GetHost(podName string, port int, challengeId string, randomId string) string {
	c := env.Get()
	return fmt.Sprintf("%s-%d-%s-%s.%s", podName, port, challengeId, randomId, c.Domain)
}
