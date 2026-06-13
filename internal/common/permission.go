package common

import "strings"

func ExtractPermissionResourceAndAction(permission string) (string, string) {
	permParts := strings.Split(permission, "#")
	return permParts[0], permParts[1]
}
