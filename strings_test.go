package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
)

func TestStringsGenerateRandom(t *testing.T) {
	fmt.Println(d.StringsGenerateRandom(25))
}

func TestBcryptPassword(t *testing.T) {
	str := "123456"
	pw, err := d.BcryptPassword(str)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pw)
}

func TestStringsToSnake(t *testing.T) {
	fmt.Println(d.StringsToSnake(""))
	fmt.Println(d.StringsToSnake("Test"))
	fmt.Println(d.StringsToSnake("TeST"))
	fmt.Println(d.StringsToSnake("test"))
}