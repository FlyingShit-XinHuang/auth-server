package client

import (
	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/pkg/encoding"
	httptransport "whispir/auth-server/pkg/transport/http"
)

const CreateClientPath = "/clients"

type httpTransport struct {
	svc Service
}

func NewHTTPTransport(svc Service) httptransport.HTTPTransport {
	return &httpTransport{
		svc: svc,
	}
}

func (h httpTransport) GetRouters() []httptransport.Route {
	return []httptransport.Route{
		httptransport.Route{
			Path:       CreateClientPath,
			Methods:    []string{"POST"},
			Endpoint:   CreateClientEndpoint(h.svc),
			ReqDecoder: encoding.DecodeRequestJsonInto(&v1alpha1.Client{}),
		},
	}
}
