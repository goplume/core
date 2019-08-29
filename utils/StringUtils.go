package utils

func StringIsNotBlank(val string) bool {
	return !StringIsBlank(val)
}

func StringIsBlank(val string) bool {
	if len(val) == 0 {
		return true
	}

	for _, char := range val {
		if IsWhiteSpaceChar(char) == false {
			return false
		}
	}

	return true
}

// contract
func StringIsEmpty(val string) bool {
	return val == ""
}

// contract
func StringIsNotEmpty(val string) bool {
	return val == ""
}

func IsWhiteSpaceChar(char int32) bool {
	return char == ' ' || char == '\n' || char == '\t' || char == '\r'
}
