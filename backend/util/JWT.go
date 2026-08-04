package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

var key = []byte("1145141919810")

// GenerateToken 生成token
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	now := time.Now()
	userClaim := &UserClaims{
		Name:     name,
		Identity: identity,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(10 * time.Minute).Unix(),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
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
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return key, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if userClaim.ExpiresAt == 0 {
		err := errors.New("token expiration is missing")
		fmt.Println(err)
		return nil, err
	}
	if !token.Valid {
		fmt.Println("token is not valid")
		return nil, errors.New("token is not valid")
	}
	return userClaim, nil
}
