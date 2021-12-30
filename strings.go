package d

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

// Deprecated: 使用StringsGenerateRandom函数替代
func GenerateRandomString(l int) string {
	return StringsGenerateRandom(l)
}

// Encrypt user password
// 加密用户密码
// https://pkg.go.dev/golang.org/x/crypto/bcrypt
func BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify user password
// 验证用户密码
func VerifyPassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, err
}

// Generate random strings
// 生成随机字符串
func StringsGenerateRandom(l int) string {
	var code = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/~!@#$%^&*()_="

	data := make([]byte, l)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < l; i++ {
		idx := rand.Intn(len(code))
		data[i] = code[idx]
	}
	return string(data)
}

// 转化为蛇形字符串，例如 XxYy to xx_yy, XxYY to xx_yy
func StringsToSnake(s string) string {
	tmp := make([]byte, 0)
	for i := 0; i < len(s); i++ {
		if i > 0 && s[i] >= 'A' && s[i] <= 'Z' {
			tmp = append(tmp, '_')
		}
		tmp = append(tmp, s[i])
	}
	return strings.ToLower(string(tmp[:]))
}