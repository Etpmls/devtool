package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
)

func TestGetField(t *testing.T) {
	d.SetField("tmp","123")
	s, err := d.GetField("tmp")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("设置tmp字段名为：123，获取tmp的值为：" + s)
	s2, err := d.GetField("no-set")
	fmt.Println("未设置no-set的字段名，获取no-set的值为：" + s2)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGetPath(t *testing.T) {
	d.SetPath("tmp","/123")
	s, err := d.GetPath("tmp")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("设置tmp路径为：/123，获取tmp的路径为：" + s)
	s2, err := d.GetPath("no-set")
	fmt.Println("未设置no-set的路径，获取no-set的路径为：" + s2)
	if err != nil {
		fmt.Println(err.Error())
	}
}
