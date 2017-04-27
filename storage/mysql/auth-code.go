package mysql

import (
	"github.com/RangelReale/osin"

	"github.com/dgrijalva/jwt-go"

	whispirjwt "whispir/auth-server/pkg/jwt"

	"fmt"
	"time"
	"whispir/auth-server/pkg/api/v1alpha1"
)

// nothing to be saved because info are stored in JWT
func (s *mysqlStorage) SaveAuthorize(*osin.AuthorizeData) error {
	return nil
}

func (s *mysqlStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	claims, err := whispirjwt.ParseAuthCode(code, func(c *whispirjwt.AuthCodeClaims) (interface{}, error) {
		user, err := s.GetUserById(c.UserId)
		if nil != err {
			fmt.Println("get user error:", err)
			return nil, fmt.Errorf("failed to get user:", err)
		}
		return []byte(user.Password), nil
	})
	if nil != err {
		fmt.Println("parse auth code err:", err)
		return nil, fmt.Errorf("failed to parse authorization code:", err)
	}

	client, err := s.GetClient(claims.ClientId)
	if nil != err {
		fmt.Println("get client error:", err)
		return nil, fmt.Errorf("failed to get client:", err)
	}

	data := &osin.AuthorizeData{
		Client:      client,
		Code:        code,
		ExpiresIn:   int32(claims.Expire - jwt.TimeFunc().Unix()),
		Scope:       claims.Scope,
		State:       claims.State,
		CreatedAt:   time.Unix(claims.IssueAt, 0),
		UserData:    &v1alpha1.User{Id: claims.UserId},
		RedirectUri: client.GetRedirectUri(),
	}
	if claims.RedirectURI != "" {
		data.RedirectUri = claims.RedirectURI
	}
	return data, nil
}

func (s *mysqlStorage) RemoveAuthorize(code string) error {
	return nil
}
