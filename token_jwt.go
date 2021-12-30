package d

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"strconv"
)

var (
	Token token
)

type token struct {
	Optional optionalToken
	enable bool
}

type optionalToken struct {
	SigningKey string
}

func (this *token) Init()  {
	// 如果没有加密签名，默认自动生成25位长度作为加密签名
	if this.Optional.SigningKey == "" {
		this.Optional.SigningKey = StringsGenerateRandom(25)
	}
	this.enable = true
}

// 获取启动的状态
func (this *token) GetEnabledStatus() bool {
	return this.enable
}

// 创建HS256标准JWT Token
func (this *token) Create(claims *jwt.StandardClaims) (string, error) {
	// https://pkg.go.dev/github.com/golang-jwt/jwt#example-NewWithClaims-StandardClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(this.Optional.SigningKey))
	if err != nil {
		return "", err
	}
	return str, nil
}

// 解析标准JWT Token
func (this *token) Parse(tokenString string) (*jwt.Token, error) {
	// https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-ErrorChecking
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(this.Optional.SigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("the token format is wrong")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, errors.New("the token has expired or has not been activated yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, errors.New("invalid token")
}

// 从token中获取用户名
func (this *token) GetSubjectByToken(tokenString string) (issuer string, err error) {
	tk, err := this.Parse(tokenString)
	if err != nil {
		return "", err
	}

	// https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		issuer, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("unable to get the iss from the token")
		}
		return issuer, nil
	} else {
		return "", errors.New("the current token is invalid")
	}
}

// 从token中获取用户ID
func (this *token) GetJwtIdByToken(tokenString string) (userId int, err error) {
	tk, err := this.Parse(tokenString)
	if err != nil {
		return 0, err
	}

	// https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac
	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		id, ok := claims["jti"].(string)
		if !ok {
			return 0, errors.New("unable to get the jti from the token")
		}

		userId, err := strconv.Atoi(id)
		if err != nil {
			return 0, errors.New("the ID in the token is not a number")
		}

		return userId, nil
	} else {
		return 0, errors.New("the current token is invalid")
	}
}