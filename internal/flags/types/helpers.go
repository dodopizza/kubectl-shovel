package types

func ContainsItemString(strings []string, item string) bool {
	for _, str := range strings {
		if str == item {
			return true
		}
	}
	return false
}
