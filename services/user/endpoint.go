package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"whispir/auth-server/pkg/api/v1alpha1"
)

func CreateUserEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		return nil, svc.CreateUser(request.(*v1alpha1.User))
	}
}
