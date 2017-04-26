package jwt_test

import (
	"errors"
	"github.com/RangelReale/osin"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"whispir/auth-server/pkg/api/v1alpha1"
	whispirJWT "whispir/auth-server/pkg/jwt"
	//"encoding/json"
	"encoding/base64"
	"strings"
)

func TestAuthCodeClaims(t *testing.T) {
	gen := whispirJWT.NewAuthCodeGenerator()
	user := &v1alpha1.User{
		Id:       123456,
		Name:     "habor",
		Password: "habor",
	}
	data := &osin.AuthorizeData{
		ExpiresIn: 3600,
		Client:    testClient{},
		UserData:  user,
	}
	token, err := gen.GenerateAuthorizeToken(data)
	t.Log(token, err)
	logEncodedClaims(t, token)

	claims := whispirJWT.NewAuthCodeClaims()
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("mismatched method")
		}
		return []byte(user.Password), nil
	})
	if nil != err {
		t.Fatal(err)
	}

	claims, ok := tok.Claims.(*whispirJWT.AuthCodeClaims)
	if !ok {
		t.Fatal("cannot convert to AuthCodeClaims")
	}

	t.Log(claims.ClientId, claims.UserId)
}

func TestAccessTokenClaims(t *testing.T) {
	gen := whispirJWT.NewAccessTokenGenerator()
	client := testClient{}
	data := &osin.AccessData{
		ExpiresIn: 3600,
		Client:    client,
	}
	token, _, err := gen.GenerateAccessToken(data, false)
	t.Log(token, err)
	logEncodedClaims(t, token)

	claims := whispirJWT.NewAccessTokenClaims()
	tok, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("mismatched method")
		}
		return []byte(client.GetSecret()), nil
	})
	if nil != err {
		t.Fatal(err)
	}

	claims, ok := tok.Claims.(*whispirJWT.AccessTokenClaims)
	if !ok {
		t.Fatal("cannot convert to AccessTokenClaims")
	}
	t.Log(claims.ClientId)
}

func logEncodedClaims(t *testing.T, token string) {
	s := strings.Split(token, ".")
	decoder := base64.RawURLEncoding
	header, _ := decoder.DecodeString(s[0])
	claims, _ := decoder.DecodeString(s[1])
	t.Log(string(header), string(claims))
}

type testClient struct {
}

func (testClient) GetId() string {
	return "foo"
}

func (testClient) GetSecret() string {
	return "foo"
}

func (testClient) GetRedirectUri() string {
	return ""
}

func (testClient) GetUserData() interface{} {
	return ""
}
