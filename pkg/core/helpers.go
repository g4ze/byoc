package core

import "strings"

func generateName(UserName string, Image, suffix string) string {
	if UserName == "" {
		return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + suffix
	}
	return UserName + "-" + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-") + "-" + suffix
}
func generateNameFromImage(Image string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(Image, "/", "-"), ".", "-"), ":", "-")
}
