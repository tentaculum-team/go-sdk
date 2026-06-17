package utils

func StartsWithInvalidNameChar(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[0] == '-' || name[0] == '\''
}

func EndsWithInvalidNameChar(name string) bool {
	if len(name) == 0 {
		return false
	}
	return name[len(name)-1] == '-' || name[len(name)-1] == '\''
}
