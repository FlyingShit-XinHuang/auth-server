package user

import (
	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/pkg/encoding"
	httptransport "whispir/auth-server/pkg/transport/http"
)

type httpTransport struct {
	svc Service
}

const CreateUserPath = "/users"

func NewHTTPTransport(svc Service) httptransport.HTTPTransport {
	return &httpTransport{
		svc: svc,
	}
}

func (h httpTransport) GetRouters() []httptransport.Route {
	return []httptransport.Route{
		httptransport.Route{
			Path:       CreateUserPath,
			Methods:    []string{"POST"},
			Endpoint:   CreateUserEndpoint(h.svc),
			ReqDecoder: encoding.DecodeRequestJsonInto(&v1alpha1.User{}),
		},
	}
}
