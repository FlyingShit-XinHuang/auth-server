package client

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"whispir/auth-server/pkg/api/v1alpha1"
)

func CreateClientEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		client := request.(*v1alpha1.Client)
		return svc.CreateClient(client)
	}
}
