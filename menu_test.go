package d_test

import (
	"encoding/json"
	"fmt"
	d "github.com/etpmls/devtool"
	"testing"
)

func TestGetMenuPath(t *testing.T) {
	p := d.GetMenuPath()
	fmt.Println("Menu Path:" + p)
	return
}

func TestMenuCreate(t *testing.T) {
	d.SetMenuPath("test/")
	type Menu struct {
		Id int `json:"id"`
		Name string `json:"name"`
	}
	var m = Menu{
		Id: 100,
		Name: "TestMenu",
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = d.MenuCreate(string(b))
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestMenuGet(t *testing.T) {
	b ,err := d.MenuGet()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println("Menu:" + string(b))
}