package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	jwt.StandardClaims
}

var key = []byte("1145141919810")

// GenerateToken 生成token
func GenerateToken(identity, name string) (string, error) {
	UserClaim := &UserClaims{
		"test",
		"e47c08df-284e-4473-b16c-70f2f56ac3f9",
		jwt.StandardClaims{},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	signedToken, err := unsignedToken.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return signedToken, nil
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImlkZW50aXR5IjoiZTQ3YzA4ZGYtMjg0ZS00NDczLWIxNmMtNzBmMmY1NmFjM2Y5In0.Ba05JtWeBRexzFDmvtuv15mPTtVXf7JO2YESqV1wieU
}

// AnalyseToken 解析
func AnalyseToken(signedToken string) (*UserClaims, error) {
	userClaim := &UserClaims{}
	token, err := jwt.ParseWithClaims(signedToken, userClaim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if !token.Valid {
		fmt.Println("token is not valid")
		return nil, err
	}
	return userClaim, nil
}
