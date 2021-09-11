package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
)

func TestCopyStructValue(t *testing.T) {
	type Source struct {
		A string
		B int
		C uint
		D bool
	}
	type Target struct {
		A string
		B int
		D bool
	}
	var source = Source{
		A: "this is a",
		B: 99,
		C: 88,
		D: true,
	}
	var target Target
	fmt.Println(source, target)
	d.CopyStructValue(source, &target)
	fmt.Println(source, target)
}