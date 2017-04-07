package auth

import (
	"github.com/RangelReale/osin"

	"context"
	"net/http"
	httptransport "whispir/auth-server/pkg/transport/http"
)

const (
	TokenPath = "/token"
	AuthPath  = "/auth"
)

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
			Path:       TokenPath,
			Endpoint:   GetAccessTokenEndpoint(h.svc),
			ResEncoder: osinEnc,
		},
		httptransport.Route{
			Path:       "/info",
			Endpoint:   InfoEndpoint(h.svc),
			ResEncoder: osinEnc,
		},
		httptransport.Route{
			Path:       AuthPath,
			Endpoint:   GetAuthCodeEndpoint(h.svc),
			ResEncoder: encodeAuthCodeResp,
		},
	}
}

// encode with osin package
func osinEnc(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*authResponse)
	defer resp.osinResp.Close()
	return osin.OutputJSON(resp.osinResp, w, resp.req)
}

func encodeAuthCodeResp(c context.Context, w http.ResponseWriter, response interface{}) error {
	switch response.(type) {
	case *bytesResponse:
		resp := response.(*bytesResponse)
		_, err := w.Write(resp.data)
		return err
	default:
		return osinEnc(c, w, response)
	}
}
