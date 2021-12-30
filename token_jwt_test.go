package d_test

import (
	"fmt"
	d "github.com/etpmls/devtool"
	"github.com/golang-jwt/jwt"
	"strconv"
	"testing"
)

func TestToken(t *testing.T) {
	d.Token.Init()
	fmt.Println("SigningKey: " + d.Token.Optional.SigningKey)

	// 创建token
	str, err := d.Token.Create(&jwt.StandardClaims{
		Audience:  "xxx",
		ExpiresAt: 0,
		Id:        "123",
		IssuedAt:  0,
		Issuer:    "",
		NotBefore: 0,
		Subject:   "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Token: " + str)

	// 解析token
	j, err := d.Token.Parse(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(j.Valid)

	// 获取用户名
	iss, err := d.Token.GetSubjectByToken(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Username: " + iss)

	// 获取用户ID
	id, err := d.Token.GetJwtIdByToken(str)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ID: " + strconv.Itoa(id))
}
