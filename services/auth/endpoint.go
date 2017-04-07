package auth

import (
	"github.com/go-kit/kit/endpoint"

	"context"
	"net/http"

	"github.com/RangelReale/osin"
)

type authResponse struct {
	osinResp *osin.Response
	req      *http.Request
}

type bytesResponse struct {
	data []byte
	req  *http.Request
}

func GetAccessTokenEndpoint(svc Service) endpoint.Endpoint {
	return endpointFactory(svc.GetAccessToken)
}

func InfoEndpoint(svc Service) endpoint.Endpoint {
	return endpointFactory(svc.Info)
}

func GetAuthCodeEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(*http.Request)
		resp, authPage := svc.GetAuthCode(req)
		if nil != resp {
			return &authResponse{
				osinResp: resp,
				req:      req,
			}, nil
		}
		return &bytesResponse{
			data: authPage,
			req:  req,
		}, nil
	}
}

func endpointFactory(handler func(req *http.Request) *osin.Response) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*http.Request)
		resp := handler(req)
		return &authResponse{
			osinResp: resp,
			req:      req,
		}, nil
	}
}
