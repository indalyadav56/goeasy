package utils

import (
	"strings"
)

func ToCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		} else {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(string(parts[i][0])) + strings.ToLower(parts[i][1:])
			}
		}
	}
	return strings.Join(parts, "")
}
