package d_test

import (
	d "github.com/Etpmls/devtool"
	"testing"
)

var (
	Host = "127.0.0.1"
	Port = "3306"
	DBName = "test"
	Username = "root"
	Password = "123456"
	TablePrefix = "ts_"
)

func TestDatabase(t *testing.T)  {
	d.Database.Host = Host
	d.Database.Port = Port
	d.Database.DBName = DBName
	d.Database.Username = Username
	d.Database.Password = Password
	d.Database.Optional.TablePrefix = TablePrefix

	type testField struct {
		TestString string
		TestString2 string
		TestInt int
		TestUint uint
		TestByte byte
		testLow string
	}
	d.Database.Optional.AutoMigrate = []interface{}{
		&testField{},
	}
	err := d.Database.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = d.Database.Migration()
	if err != nil {
		t.Fatal(err)
	}
}