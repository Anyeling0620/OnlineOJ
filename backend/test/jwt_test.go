package test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	jwt.StandardClaims
}

var key = []byte("1145141919810")

func TestGenerateToken(t *testing.T) {
	UserClaim := &UserClaims{
		"test2",
		"c4d038b4bed09fdb1471ef51ec3a32cd",
		jwt.StandardClaims{},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	signedToken, err := unsignedToken.SignedString([]byte(key))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signedToken)
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImlkZW50aXR5IjoiZTQ3YzA4ZGYtMjg0ZS00NDczLWIxNmMtNzBmMmY1NmFjM2Y5In0.Ba05JtWeBRexzFDmvtuv15mPTtVXf7JO2YESqV1wieU
}

func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImlkZW50aXR5IjoiZTQ3YzA4ZGYtMjg0ZS00NDczLWIxNmMtNzBmMmY1NmFjM2Y5In0.Ba05JtWeBRexzFDmvtuv15mPTtVXf7JO2YESqV1wieU"
	userClaim := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.Valid {
		fmt.Println(userClaim)
	}
}
