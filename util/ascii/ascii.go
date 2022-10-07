package ascii

func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func ToLower(c byte) byte {
	if IsUpper(c) {
		return c + 0x20
	}
	return c
}

func ToUpper(c byte) byte {
	if IsLower(c) {
		return c - 0x20
	}
	return c
}
