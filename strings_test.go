package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	fmt.Println(d.GenerateRandomString(25))
}

func TestBcryptPassword(t *testing.T) {
	str := "123456"
	pw, err := d.BcryptPassword(str)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pw)
}
