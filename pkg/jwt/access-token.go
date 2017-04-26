package jwt

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"whispir/auth-server/pkg/api/v1alpha1"
)

type AccessTokenClaims struct {
	*CommonClaims
}

func NewAccessTokenGenerator() *AccessTokenClaims {
	return &AccessTokenClaims{}
}

func NewAccessTokenClaims() *AccessTokenClaims {
	return &AccessTokenClaims{}
}

// implement interface jwt.Claims
func (a *AccessTokenClaims) Valid() error {
	return NewValidator(a).Valid()
}

// implement interface osin.AccessTokenGen
func (a *AccessTokenClaims) GenerateAccessToken(data *osin.AccessData, generaterefresh bool) (accesstoken string, refreshtoken string, err error) {
	now := jwt.TimeFunc().Unix()
	user, ok := data.UserData.(*v1alpha1.User)

	a = &AccessTokenClaims{
		&CommonClaims{
			Expire:    now + int64(data.ExpiresIn),
			NotBefore: now,
			IssueAt:   now,
			Scope:     data.Scope,
			ClientId:  data.Client.GetId(),
		},
	}
	if ok {
		a.UserId = user.Id
	}

	token := jwt.NewWithClaims(signMethod, a)

	tokstr, err := token.SignedString([]byte(data.Client.GetSecret()))

	return tokstr, "", nil
}

func ParseAccessToken(token string, keygen func(*AccessTokenClaims) (interface{}, error)) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	tok, err := jwt.ParseWithClaims(token, claims, KeyGenerator(claims, func(obj interface{}) (interface{}, error) {
		c := obj.(*AccessTokenClaims)
		return keygen(c)
	}))
	if nil != err {
		fmt.Println("parse jwt err:", err)
		return nil, fmt.Errorf("failed to parse jwt:", err)
	}

	claims = tok.Claims.(*AccessTokenClaims)
	return claims, nil
}
