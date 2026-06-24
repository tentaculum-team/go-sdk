package utils

func IsNumericOnly(username string) bool {
	if len(username) == 0 {
		return false
	}

	for i := 0; i < len(username); i++ {
		if username[i] < '0' || username[i] > '9' {
			return false
		}
	}

	return true
}
