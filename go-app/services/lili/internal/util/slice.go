package util

func Contains[T comparable](slice []T, target T) bool {
	for _, elm := range slice {
		if target == elm {
			return true
		}
	}
	return false
}
