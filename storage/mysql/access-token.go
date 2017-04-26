package mysql

import (
	"github.com/RangelReale/osin"

	"fmt"
	"time"

	whispirjwt "whispir/auth-server/pkg/jwt"

	"github.com/dgrijalva/jwt-go"
	"whispir/auth-server/pkg/api/v1alpha1"
)

// nothing to be saved because info are stored in JWT
func (s *mysqlStorage) SaveAccess(*osin.AccessData) error {
	return nil
}

func (s *mysqlStorage) LoadAccess(token string) (*osin.AccessData, error) {
	var client osin.Client
	claims, err := whispirjwt.ParseAccessToken(token, func(c *whispirjwt.AccessTokenClaims) (interface{}, error) {
		var err error
		client, err = s.GetClient(c.ClientId)
		if nil != err {
			fmt.Println("get client error:", err)
			return nil, fmt.Errorf("failed to get client: %v\n", err)
		}

		return []byte(client.GetSecret()), nil
	})
	if nil != err {
		fmt.Println("parse access token err:", err)
		return nil, fmt.Errorf("failed to parse access token: %v\n", err)
	}

	data := &osin.AccessData{
		Client:      client,
		AccessToken: token,
		ExpiresIn:   int32(claims.Expire - jwt.TimeFunc().Unix()),
		Scope:       claims.Scope,
		CreatedAt:   time.Unix(claims.IssueAt, 0),
	}
	if claims.UserId > 0 {
		data.UserData = &v1alpha1.User{Id:claims.UserId}
	}

	return data, nil
}

func (s *mysqlStorage) RemoveAccess(token string) error {
	return nil
}

func (s *mysqlStorage) LoadRefresh(token string) (*osin.AccessData, error) {
	return nil, fmt.Errorf("not supported")
}

func (s *mysqlStorage) RemoveRefresh(token string) error {
	return nil
}
