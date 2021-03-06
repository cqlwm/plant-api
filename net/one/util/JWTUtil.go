package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	secretKey = "cq.lwm.test.plan.2dgr3ij6alva8.key" //私钥
)

//自定义Claims
type CustomClaims struct {
	UserId int
	jwt.StandardClaims
}

// 生成Token
func BuildToken(userId int) string {
	//生成token
	maxAge := 60 * 60 * 24
	customClaims := &CustomClaims{
		UserId: userId, //用户id
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(maxAge) * time.Second).Unix(), // 过期时间，必须设置
		},
	}

	//采用HMAC SHA256加密算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("token: %v\n", tokenString)
	return tokenString
}

// 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// 从token中拿到UserId
func UserInfo(c *gin.Context) (*CustomClaims, error) {
	token := c.GetHeader("token")
	if len(token) == 0 {
		// 没有Token参数
		log.Println("没有Token参数", token)
		return nil, errors.New("没有Token参数")
	}

	claims, err := ParseToken(token)
	if err != nil {
		// Token无效
		log.Println("没有Token参数", err)
		return nil, err
	}

	// 返回CustomClaims
	return claims, nil

}

func TokenTest() {
	token := BuildToken(1234567)

	//解析token
	ret, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("userinfo: %v\n", ret.UserId)
}
