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

Currently We Support: 

Library: **Log**, **Config**, **Customize**, **Database**, **Validator**, **Token**, **Strings**, **Stuct**, **Cache**, **Menu**

Module: **Recaptcha**, **Gin**, **Dingtalk**

# Library

We recommend initializing the Log library first. In order to initialize other libraries, you can easily record the errors of other libraries.

## Log

**Introduce**

We use ***sirupsen/logrus*** package. [sirupsen/logrus: Structured, pluggable logging for Go. (github.com)](https://github.com/sirupsen/logrus)

Through the log library, you can output log information faster.



**Parameter**

> Optional.Level ( logrus.Level | optional )

Log level defined by logrus, such as TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel.

*For more information, please read logrus related documents*



> Optional.Filename ( string | optional )

The file name(include path) of the log storage, if it is empty, it will be stored in log/app.log by default.

*For more information, please read logrus related documents*



> Optional.Maxsize ( int| optional )
>
> Optional.MaxBackups ( int| optional )
>
> Optional.MaxAge ( int| optional )
>
> Optional.Compress ( bool| optional )

*For more information, please read logrus related documents*



**Example**

> main.go

```go
func init()  {
	d.Log.Optional.Compress = true
	d.Log.Init()
}

func main() {
    logrus.Info("This is an example.")
}
```



## Config

**Introduce**

We use ***go-yaml/yaml*** package. [go-yaml/yaml: YAML support for the Go language. (github.com)](https://github.com/go-yaml/yaml)

Through the *Config* library, you can save the application configuration in the specified file, and read them through the library to achieve flexible reading of variables.



**Parameter**

> FilePath ( string | required )
> 

Specifies the directory for storing configuration files

Example:
```go
d.Config.FilePath = "storage/config/config.yaml"
```

> ConfigAddr ( interface{} | required | pointer )
>

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



**Example**

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
func init() {
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
        log.Fatal(err)
        return
    }
}

func main(){
    fmt.Println(Config.App.HttpPort)
}
```



## Database

**Introduce**

We use ***go-gorm/gorm*** package. [go-gorm/gorm: The fantastic ORM library for Golang, aims to be developer friendly (github.com)](https://github.com/go-gorm/gorm)

The default database of the ***Database*** library is MySQL. Through the ***Database*** library, you can quickly initialize MySQL.



**Parameter**

> Host ( string | required )

MySQL server host

> Port ( string | required )

MySQL server port

> DBName ( string | required )

MySQL database name

> Username ( string | required )

Database username

> Password ( string | required )

Database password

> Optional.TablePrefix ( string | optional )

Database table prefix 

> Optional.AutoMigrate ( []interface{} | optional )

Quickly create a table structure. Refer to the gorm document for details.

> Optional.Charset ( string | optional )

Database Charset. Default utf8mb4.

> Optional.DoNotParseTime ( bool| optional )

Default false. Whether the database is parsed time.

> Optional.Location ( string | optional )

Default "Local". The location corresponding to the database

> Optional.FuzzySearch ( string | optional )

Define keywords for fuzzy search. Modify the definition only when changing the database. The default is "LIKE" in MySQL.



**Function**

> Migration() error

Create a table structure. If you do not use this method, the table structure will not be created automatically by default.



**Example**

```go
func init(){
    type Test struct{
        A string
        B string
    }
    
    d.Database.Host = "127.0.0.1"
    d.Database.Port = "3306"
    d.Database.DBName = "dbname"
	d.Database.Username = "root"
    d.Database.Password = "root"
    d.Database.Optional.TablePrefix = "sql_"
	d.Database.Optional.AutoMigrate = []interface{}{
		&Test{},
	}
	err = d.Database.Init()
	if err != nil {
		logrus.Fatal(err)
	}
	err = d.Database.Migration()
	if err != nil {
		logrus.Fatal(err)
	}
}
```





## Validator

**Introduce**

We use ***go-playground/validator*** package. [go-playground/validator: :100:Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving (github.com)](https://github.com/go-playground/validator)

Package validator implements value validations for structs and individual fields based on tags.



**Parameter**

> Optional.OverrideInit ( func() error | optional )

If it is empty, the ***ValidatorClient*** is initialized as the default validator object for ***go-playground/validator*** by default. If you need more advance features, such as a multi-language validator, you can define this method.



**Example**

```go
import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Translator ut.Translator
)

func init(){
    	d.Validator.Optional.OverrideInit = func() error {
		var (
			uni      *ut.UniversalTranslator
		)
		zh := zh.New()
		uni = ut.New(zh, zh)
		Translator, _ = uni.GetTranslator("zh")
		return zh_translations.RegisterDefaultTranslations(d.ValidatorClient, Translator)
	}
    
    err := d.Validator.Init()
	if err != nil {
		logrus.Fatal(err)
		return
	}
}
```



## Token

**Introduce**

We use ***golang-jwt/jwt*** package. [golang-jwt/jwt: Community maintained clone of https://github.com/dgrijalva/jwt-go](https://github.com/golang-jwt/jwt)

Through the ***Token*** library, you can quickly generate or parse tokens.



**Parameter**

> Optional.SigningKey ( string | optional )

Set your encryption key, if it is empty, it will be generated by default.



**Function**

> Create(claims *jwt.StandardClaims) (string, error)

Create HS256 standard JWT Token



> Parse(tokenString string) (*jwt.Token, error)

Parse HS256 standard JWT Token



> GetSubjectByToken(tokenString string) (issuer string, err error)

Get the username from the token



> GetJwtIdByToken(tokenString string) (userId int, err error)

Obtain the user ID from the token



**Example**

> main.go

```go
func init(){
    d.Token.Init()
}
```



# Module

## Recaptcha

**Introduce**

Only need a simple configuration, you can use the Recaptcha method to easily verify the verification code.



**Parameter**

> Secret ( string | required )

Your recaptcha secret

> Optional.Host ( string | optional )

The default domain name of recaptcha, www.google.com, such as Mainland China, please set www.recaptcha.net.



**Function**

> VerifyCaptcha(response string) (error)

Verify the response code returned from the front end



**Example**

```go
	// 初始化recaptcha
	d_recaptcha.Recaptcha.Secret = "YOUR_RECAPTCHA_SECRET"
	d_recaptcha.Recaptcha.Optional.Host = "www.recaptcha.net"
	err = d_recaptcha.Recaptcha.Init()
	if err != nil {
		logrus.Fatal(err)
		return
	}
```



## Gin

**Get started**

Import package

```go
import "github.com/Etpmls/devtool/gin"
```

**Function**

> SetCors(router *gin.Engine, allowHeaders []string)

Set the header field name allowed by CORS

Example

```go
func main()  {
	r := gin.Default()
	d_gin.SetCors(r, []string{"token", "language"})
	r.Run()
}
```

### Validate

**Function**

> Validate(c *gin.Context, structAddr interface{}, translator ut.Translator) error

Verify that the request meets the requirements

**Example**

```go
err := d_gin.Validate(c, json, nil)
if err != nil {
    return err
}
```



### Auth

**Depend**

Database, Log, Validator Libraries

**Parameter**

> Optional.SkipInsertingInitializationData ( bool | optional )

Default false. Skip to initialize permission related data

> TokenExpirationTime ( time.Duration | optional )

Default 43200. Token expiration time, unit s (seconds)

**Function**

> GetUserIdByRequest(c *gin.Context) (uint, error)

Get User ID by request token

> GetUserByRequest(c *gin.Context) (u User, err error)

Get User information by request token

**Example**

```go
	err = d_gin.Auth.Init()
	if err != nil {
		logrus.Fatal(err)
		return
	}
```



## Dingtalk

**Parameter**

> AgentId ( string | required )
>
> AppKey ( string | required )
>
> AppSecret ( string | required )
>
> Optional.CorpId ( string | optional )

Obtained from Dingtalk

**Function**

> GetAccessToken() (s string, err error)

Get Dingding access token

> SendTextNotificationsOfWork(userID string, content string) (err error)

Send text-based job notifications

> SendOaNotificationsOfWork(userID string, m OaMessage) (err error)

Send OA type job notice

> GetUserInfo(access_token string, code string) (userInfo GetUserInfoResponse, err error)

Password-free login,Get user information

> GetUserId(access_token string, code string) (userId string, err error)

Password-free login,Get user id

> GenerateJumpableUrl(appurl string) (string, error)

The URL prefix required by the DingTalk workbench

**Example**

```go
	d_dingtalk.Dingtalk.AgentId = "YOUR_AGENTID"
	d_dingtalk.Dingtalk.AppKey = "YOUR_APPKEY"
	d_dingtalk.Dingtalk.AppSecret = "YOUR_APPSECRET"
	err = d_dingtalk.Dingtalk.Init()
	if err != nil {
		logrus.Fatal(err)
		return
	}
```



