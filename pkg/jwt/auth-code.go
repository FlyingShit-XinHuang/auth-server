package jwt

import (
	"errors"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"whispir/auth-server/pkg/api/v1alpha1"
)

func NewAuthCodeGenerator() *AuthCodeClaims {
	return &AuthCodeClaims{}
}

func NewAuthCodeClaims() *AuthCodeClaims {
	return &AuthCodeClaims{}
}

type AuthCodeClaims struct {
	*CommonClaims
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirect_url,omitempty"`
}

func (c *AuthCodeClaims) missingFields() bool {
	if c.CommonClaims.missingFields() || c.UserId == 0 {
		return true
	}
	return false
}

func (c *AuthCodeClaims) invalidScope() bool {
	return false
}

// implement interface jwt.Claims
func (c *AuthCodeClaims) Valid() error {
	if c.invalidScope() {
		fmt.Println("scope is invalid", c.Scope)
		return errors.New("Invalid scope")
	}
	return NewValidator(c).Valid()
}

// implement interface osin.AuthorizeTokenGen
func (c *AuthCodeClaims) GenerateAuthorizeToken(data *osin.AuthorizeData) (ret string, err error) {
	now := jwt.TimeFunc().Unix()
	user, ok := data.UserData.(*v1alpha1.User)
	if !ok {
		fmt.Println("failed to get a user object from", data.UserData)
		return "", errors.New("no user found")
	}

	c = &AuthCodeClaims{
		&CommonClaims{
			Expire:    now + int64(data.ExpiresIn),
			NotBefore: now,
			IssueAt:   now,
			Scope:     data.Scope,
			ClientId:  data.Client.GetId(),
			UserId:    user.Id,
		},
		data.State,
		data.RedirectUri,
	}
	token := jwt.NewWithClaims(signMethod, c)

	return token.SignedString([]byte(user.Password))
}

func ParseAuthCode(code string, keygen func(*AuthCodeClaims) (interface{}, error)) (*AuthCodeClaims, error) {
	claims := &AuthCodeClaims{}
	token, err := jwt.ParseWithClaims(code, claims, KeyGenerator(claims, func(obj interface{}) (interface{}, error) {
		c := obj.(*AuthCodeClaims)
		return keygen(c)
	}))
	if nil != err {
		fmt.Println("parse jwt err:", err)
		return nil, fmt.Errorf("failed to parse jwt: %v\n", err)
	}

	claims = token.Claims.(*AuthCodeClaims)
	return claims, nil
}
