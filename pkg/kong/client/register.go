package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	API_NAME = "KONG_API_NAME"
	API_URL  = "KONG_API_URL"
	HOST     = "KONG_HOST"
)

var (
	ApiName string
	ApiURL string
	KongHost string
)

func RegisterAPI() error {
	var ok bool
	ApiName, ok = os.LookupEnv(API_NAME)
	if !ok {
		ApiName = "auth"
	}

	ApiURL, ok = os.LookupEnv(API_URL)
	if !ok {
		ApiURL = "http://localhost:18080"
	}

	KongHost, ok = os.LookupEnv(HOST)
	if !ok {
		fmt.Println("Running without kong because the host is not specified")
		fmt.Println("It can be specified with env", HOST)
		return nil
	}

	hostURL, err := url.Parse(KongHost)
	if nil != err {
		return fmt.Errorf("failed to parse kong host:", err)
	}
	hostURL.Path = "/apis/" + ApiName

	// get api
	resp, err := http.Get(hostURL.String())
	if nil != err {
		return fmt.Errorf("failed to get api info:", err)
	}
	if http.StatusNotFound != resp.StatusCode {
		fmt.Println("api already exists")
		return nil
	}
	resp.Body.Close()

	hostURL.Path = "/apis"
	resp, err = http.PostForm(hostURL.String(), url.Values{
		"name":         {ApiName},
		"uris":         {"/" + ApiName},
		"upstream_url": {ApiURL},
	})
	defer resp.Body.Close()

	if nil != err {
		return fmt.Errorf("failed to add api:", err)
	}
	if http.StatusCreated != resp.StatusCode {
		return fmt.Errorf("failed to add api: (%d)%s", resp.StatusCode, resp.Status)
	}

	return nil
}

func GetAPIPathPrefix() string {
	if KongHost == "" {
		return ""
	}
	return "/" + ApiName
}