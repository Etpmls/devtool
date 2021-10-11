# Introduction
devtool is a development tool library. We have packaged various tools, and you can use them as needed for rapid development.

# How To Use?
**1. Import**
```go
import (
	d "github.com/Etpmls/devtool"
)
```
**2. Go mod**
```shell
go mod vendor
```
**3. Use Your Library (Such As 'Log')**
```go
package main

import (
	d "github.com/Etpmls/devtool"
	log "github.com/sirupsen/logrus"
)

func main() {
	d.Log.Compress = true
	d.Log.Init()

	log.Info("This is Info.")
}
```

Currently We Support: **Log**, **Database**, **Validator**, **Token**, **Strings**, **Config**

# Library
## Config
> FilePath
> 
> ( string | required )

Specifies the directory for storing configuration files

Example:
```go
d.Config.FilePath = "storage/config/config.yaml"
```

> ConfigAddr
>
> ( interface{} | required | pointer)

An empty structure pointer for your custom configuration. Pass pointers to facilitate binding configuration to your structure

Example:

```go
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
```

### Initialization Example
> test/config_yaml.yaml

```yaml
app:
  http-port: "8081"
  key: "123456"
  enable-database: false
  token-expiration-time: 86400
  testfield: "good"
```
> main.go
```go
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
```