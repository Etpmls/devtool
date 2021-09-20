package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"os"
	"testing"
)

func TestGetUploadPath(t *testing.T) {
	p := d.GetUploadPath()
	fmt.Println("Upload Path:" + p)
	return
}

func TestImageValidate(t *testing.T) {
	f, err := os.OpenFile("test/test_file_noimage.png", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = d.ImageValidate(f)
	fmt.Println(err)

	f, err = os.OpenFile("test/test_file_image.png", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
		return
	}
	_, err = d.ImageValidate(f)
	fmt.Println(err)
}

func TestFileCheck(t *testing.T) {
	err := d.FileCheck("test/test_file_noimage.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	err = d.FileCheck("test/notExist.png")
	if err != nil {
		t.Fatal(err)
		return
	}
	return
}

func TestFilePathValidate(t *testing.T) {
	err := d.FilePathValidate("test/test_file_image.png", []string{"test/"})
	fmt.Println("Path:test/test_file_image.png, Upload Path:test/")
	if err != nil {
		fmt.Println("Error:" + err.Error())
		t.Fatal(err)
	}
	err = d.FilePathValidate("test/test_file_image.png", []string{"./test/"})
	fmt.Println("Path:test/test_file_image.png, Upload Path:./test/")
	if err != nil {
		fmt.Println("Error:" + err.Error())
		t.Fatal(err)
	}
	err = d.FilePathValidate("./test/test_file_image.png", []string{"test/"})
	fmt.Println("Path:./test/test_file_image.png, Upload Path:test/")
	if err != nil {
		fmt.Println("Error:" + err.Error())
		t.Fatal(err)
	}
}

func TestFileDelete(t *testing.T) {
	f, err := os.Create("test/TestFileDelete.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	f.Close()

	err = d.FileDelete("test/TestFileDelete.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestFileBatchDelete(t *testing.T) {
	f1, err := os.Create("test/TestFileBatchDelete-1.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	f1.Close()
	f2, err := os.Create("test/TestFileBatchDelete-2.txt")
	if err != nil {
		t.Fatal(err)
		return
	}
	f2.Close()
	err = d.FileBatchDelete([]string{"test/TestFileBatchDelete-1.txt", "test/TestFileBatchDelete-2.txt", "test/TestFileBatchDelete-3.txt"})
	if err != nil {
		t.Fatal(err)
		return
	}
}