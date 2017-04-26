package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
)

var signMethod = jwt.SigningMethodHS256

type CommonClaims struct {
	Expire    int64  `json:"exp"`
	NotBefore int64  `json:"nbf"`
	IssueAt   int64  `json:"iat"`
	Scope     string `json:"scope,omitempty"`
	ClientId  string `json:"client_id"`
	//TODO: check the data type
	UserId int `json:"user_id,omitempty"`
	//timeout int64
}

// implement interface jwtValidatorI
func (t *CommonClaims) missingFields() bool {
	return t.Expire <= 0 || t.NotBefore <= 0 || t.IssueAt <= 0 || t.ClientId == ""
}

func (t *CommonClaims) invalidTime(now int64) bool {
	return t.Expire <= t.NotBefore || t.Expire <= t.IssueAt
}

func (t *CommonClaims) verifyNotBefore(now int64) bool {
	return now >= t.NotBefore
}

func KeyGenerator(expectedType interface{}, gen func(interface{}) (interface{}, error)) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			fmt.Println("unexpected sign method", token.Method)
			return nil, errors.New("mismatched method")
		}

		objType := reflect.TypeOf(token.Claims)
		if objType != reflect.TypeOf(expectedType) {
			fmt.Println("unexpected type:", objType)
			return nil, errors.New("unknown claims type")
		}
		return gen(token.Claims)
	}
}
