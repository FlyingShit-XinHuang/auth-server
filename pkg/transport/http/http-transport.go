package http

import (
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
	"whispir/auth-server/pkg/encoding"
)

type HTTPTransport interface {
	GetRouters() []Route
}

type Route struct {
	Path       string
	Methods    []string
	Endpoint   endpoint.Endpoint
	ReqDecoder kithttp.DecodeRequestFunc
	ResEncoder kithttp.EncodeResponseFunc
}

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	return &Router{
		Router: mux.NewRouter(),
	}
}

func (router *Router) AddHandlers(h HTTPTransport) {
	for _, r := range h.GetRouters() {
		if len(r.Methods) < 1 {
			router.Handle(r.Path, r.newServer())
		} else {
			router.Methods(r.Methods...).Path(r.Path).Handler(r.newServer())
		}
	}
}

func (r *Route) newServer() *kithttp.Server {
	// set default decoder and encoder
	if nil == r.ReqDecoder {
		r.ReqDecoder = encoding.NoopDecodeRequest
	}
	if nil == r.ResEncoder {
		r.ResEncoder = encoding.EncodeJsonResponse
	}
	return kithttp.NewServer(r.Endpoint, r.ReqDecoder, r.ResEncoder)
}
