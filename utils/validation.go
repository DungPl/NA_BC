package utils

func IsValidValueOfConstant(role string, constantValues []string) bool {
	for _, r := range constantValues {
		if r == role {
			return true
		}
	}
	return false
}
