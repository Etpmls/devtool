package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	fmt.Println(d.GenerateRandomString(25))
}
