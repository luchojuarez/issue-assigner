package utils

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsAny(arr []string, anyOfThis []string) bool {
	for _, elem := range arr {
		if Contains(anyOfThis, elem) {
			return true
		}
	}
	return false
}
