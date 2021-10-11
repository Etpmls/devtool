package d_test

import (
	"fmt"
	d "github.com/Etpmls/devtool"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	type Configuration struct {
		App struct{
			HttpPort string	`yaml:"http-port"`
			Key string
			EnableDatabase bool	`yaml:"enable-database"`
			TokenExpirationTime time.Duration	`yaml:"token-expiration-time"`
			TestField string
		}
	}

	var Config = Configuration{}

	d.Config.FilePath = "test/config_yaml.yaml"
	d.Config.ConfigAddr = &Config
	err := d.Config.Init()
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(Config)
	fmt.Println(Config.App.HttpPort)
	fmt.Println(Config.App.Key)
	fmt.Println(Config.App.EnableDatabase)
	fmt.Println(Config.App.TokenExpirationTime)
	fmt.Println(Config.App.TestField)
}