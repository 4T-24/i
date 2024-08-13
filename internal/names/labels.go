package names

func GetCommonLabels(challengeId, teamId, instanceId string) map[string]string {
	return map[string]string{
		"i.4ts.fr/challenge":           challengeId,
		"i.4ts.fr/team":                teamId,
		"i.4ts.fr/instance":            instanceId,
		"app.kubernetes.io/managed-by": "atsi",
	}
}
