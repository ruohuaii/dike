package dike

import "testing"

func Test_GenerateCheck(t *testing.T) {
	GenerateCheck(&TBuilder{
		FavCol: &Favorite{},
	}, "1.go")
}
