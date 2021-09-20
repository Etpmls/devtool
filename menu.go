package d

import (
	"io/ioutil"
	"os"
)

// 获取菜单路径，可自定义，如无设置，则获取默认路径
func GetMenuPath() string {
	// 如果没有自定义路径，将使用默认路径
	menuDir, err := GetPath(PathMenuDir)
	if err != nil {
		menuDir = "storage/menu/"
	}
	return menuDir
}

// 设置菜单文件保存路径
func SetMenuPath(path string) {
	// 如果没有自定义路径，将使用默认路径
	SetPath(PathMenuDir, path)
	return
}

// 创建Json的菜单
func MenuCreate(jsonMenu string) error {
	// Move files
	// 移动文件
	err := os.Rename(GetMenuPath() + "menu.json", GetMenuPath() + "menu.json.bak")
	if err != nil {
		return err
	}

	// Write file
	// 写入文件
	var s = []byte(jsonMenu)
	err = ioutil.WriteFile(GetMenuPath() + "menu.json", s, 0666)
	if err != nil {
		// 还原历史菜单
		os.Rename(GetMenuPath() + "menu.json.bak", GetMenuPath() + "menu.json")
		return err
	}

	return nil
}

// 从文件中获取菜单
func MenuGet() ([]byte, error) {
	ctx, err := ioutil.ReadFile("./" + GetMenuPath() + "/menu.json")
	if err != nil {
		return nil, err
	}

	return ctx, nil
}