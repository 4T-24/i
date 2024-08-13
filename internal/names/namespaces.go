package names

import (
	"fmt"
)

func GetNamespaceName(challengeId, instanceId string) string {
	return fmt.Sprintf("atsi-%s-%s", challengeId, instanceId)
}
