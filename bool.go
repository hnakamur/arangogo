package arangogo

var myTrue bool = true
var myFalse bool = false

func TruePtr() *bool {
	return &myTrue
}

func FalsePtr() *bool {
	return &myFalse
}
