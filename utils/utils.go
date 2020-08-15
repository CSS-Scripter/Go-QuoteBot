package utils

// Check checks if an error occured. If it did, it will panic
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
