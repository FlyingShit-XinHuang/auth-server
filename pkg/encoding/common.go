package encoding

import (
	"context"
	"encoding/json"
	"net/http"

	"bytes"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"io/ioutil"
)

// Json decode request body
func DecodeRequestJsonInto(expected interface{}) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		if err := json.NewDecoder(r.Body).Decode(expected); nil != err {
			return nil, err
		}
		return expected, nil
	}
}

// noop request decoder
func NoopDecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

// Json encode response body
func EncodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Json encode request body
func EncodeJsonRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Json decode response body
func DecodeResponseJsonInto(expected interface{}) kithttp.DecodeResponseFunc {
	return func(_ context.Context, response *http.Response) (interface{}, error) {
		if http.StatusOK != response.StatusCode {
			return nil, fmt.Errorf("Response Error", response.Status)
		}
		if err := json.NewDecoder(response.Body).Decode(expected); nil != err {
			return nil, fmt.Errorf("Failed to decode response: %v", err)
		}
		return expected, nil
	}
}
