package dike

import "testing"

func Test_Work(t *testing.T) {
	Work(&TBuilder{
		FavCol: &Favorite{},
	}, "1.go")
}
