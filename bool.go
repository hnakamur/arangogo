package arangogo

func TruePtr() *bool {
	v := true
	return &v
}

func FalsePtr() *bool {
	v := false
	return &v
}
