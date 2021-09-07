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

Currently We Support: **Log**, **Database**, **Validator**, **Token**, **Strings**