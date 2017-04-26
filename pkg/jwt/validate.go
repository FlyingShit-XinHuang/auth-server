package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type jwtValidatorI interface {
	missingFields() bool
	invalidTime(now int64) bool
	verifyNotBefore(now int64) bool
}

type jwtValidator struct {
	jwtValidatorI
}

func NewValidator(i jwtValidatorI) *jwtValidator {
	return &jwtValidator{i}
}

// implement interface jwt.Claims
func (v *jwtValidator) Valid() error {
	if v.missingFields() {
		fmt.Printf("missing some fields: %#v\n", *v)
		return errors.New("Missing fields")
	}
	now := jwt.TimeFunc().Unix()
	if v.invalidTime(now) {
		fmt.Println("there are invalid time fields", *v)
		return errors.New("Invalid time exists")
	}
	if !v.verifyNotBefore(now) {
		fmt.Println("use code before allowed time", *v)
		return errors.New("Code should be used later")
	}

	return nil
}

//type nilValidator struct {}
//
//func (nilValidator) missingFields() bool {
//	return true
//}
//
//func (nilValidator) invalidTime(now int64) bool {
//	return true
//}
//
//func (nilValidator) verifyNotBefore(now int64) bool {
//	return false
//}
//
//func (nilValidator) verifyExpire(now int64) bool {
//	return false
//}
//
//func (nilValidator) invalidScope() bool {
//	return true
//}
