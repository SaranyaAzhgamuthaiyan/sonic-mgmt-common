package translib

import "strings"

func isGuageNode(key string, nodes []string) bool {
	for _, node := range nodes {
		if strings.Contains(key, node) {
			return true
		}
	}
	return false
}