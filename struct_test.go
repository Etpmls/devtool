package d_test

import (
	"fmt"
	d "github.com/etpmls/devtool"
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

func TestStructToMap(t *testing.T) {
	var a = struct {
		AaAa string
		BbBB int
		Cccc bool
	}{
		AaAa: "test",
		BbBB: 999,
		Cccc: true,
	}
	m1, err := d.StructToMap(a)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(m1)
	m2, err := d.StructToMap(&a)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(m2)
}

func TestStructToMapSnakeKey(t *testing.T) {
	var a = struct {
		AaAa string
		BbBB int
		Cccc bool
	}{
		AaAa: "test",
		BbBB: 999,
		Cccc: true,
	}
	m1, err := d.StructToMapSnakeKey(a)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(m1)
	m2, err := d.StructToMapSnakeKey(&a)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(m2)
}