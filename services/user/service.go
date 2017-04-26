package user

import (
	"context"
	"net/http"
	"net/url"

	httptransport "github.com/go-kit/kit/transport/http"

	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/pkg/encoding"
	"whispir/auth-server/storage"
)

type Service interface {
	CreateUser(*v1alpha1.User) error
}

type userService struct {
	storage storage.OAuth2Storage
}

func NewBasiceService(storage storage.OAuth2Storage) Service {
	return &userService{
		storage,
	}
}

func (u *userService) CreateUser(user *v1alpha1.User) error {
	return u.storage.CreateUser(user)
}

type serviceForHTTPClient struct {
	url *url.URL
}

func NewServiceForHTTPClient(url *url.URL) Service {
	return &serviceForHTTPClient{
		url: url,
	}
}

func (s *serviceForHTTPClient) CreateUser(user *v1alpha1.User) error {
	_, err := httptransport.NewClient(
		http.MethodPost,
		s.url,
		encoding.EncodeJsonRequest,
		encoding.DecodeResponseJsonInto(&v1alpha1.User{}),
	).Endpoint()(context.Background(), user)

	return err
}
