package client

import (
	"github.com/pborman/uuid"

	httptransport "github.com/go-kit/kit/transport/http"

	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"context"
	"fmt"
	"net/http"
	"net/url"
	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/pkg/encoding"
	"whispir/auth-server/storage"
)

type Service interface {
	CreateClient(*v1alpha1.Client) (newclient *v1alpha1.Client, err error)
}

type clientService struct {
	storage storage.OAuth2Storage
}

func NewBasicService(storage storage.OAuth2Storage) Service {
	return &clientService{
		storage,
	}
}

func (c *clientService) CreateClient(client *v1alpha1.Client) (*v1alpha1.Client, error) {
	secret := make([]byte, 16)
	rand.Read(secret)
	id := base64.RawURLEncoding.EncodeToString([]byte(uuid.NewRandom()))
	secretHex := hex.EncodeToString(secret)
	newClient := &v1alpha1.Client{
		Id:          id,
		Secret:      secretHex,
		RedirectURL: client.RedirectURL,
		Name:        client.Name,
	}
	if err := c.storage.CreateClient(newClient); nil != err {
		return nil, err
	}
	return newClient, nil
}

type clientServiceForHTTPClient struct {
	url *url.URL
}

func NewServiceForHTTPClient(url *url.URL) Service {
	return &clientServiceForHTTPClient{
		url: url,
	}
}

func (c *clientServiceForHTTPClient) CreateClient(client *v1alpha1.Client) (*v1alpha1.Client, error) {
	result, err := httptransport.NewClient(
		http.MethodPost,
		c.url,
		encoding.EncodeJsonRequest,
		encoding.DecodeResponseJsonInto(&v1alpha1.Client{}),
	).Endpoint()(context.Background(), client)

	if err != nil {
		return nil, err
	}

	newclient, ok := result.(*v1alpha1.Client)
	if !ok {
		return nil, fmt.Errorf("Unexpected type: %v", result)
	}

	return newclient, nil
}
