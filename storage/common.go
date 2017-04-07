package storage

import (
	"github.com/RangelReale/osin"
	"whispir/auth-server/pkg/api/v1alpha1"
)

type OAuth2Storage interface {
	osin.Storage
	CreateClient(client *v1alpha1.Client) error
	CreateUser(user *v1alpha1.User) error
	GetUserByNameAndPassword(name, password string) (*v1alpha1.User, error)
}
