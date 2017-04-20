package commands

import (
	//"flag"
	"fmt"
	"github.com/RangelReale/osincli"
	"github.com/spf13/cobra"
	"net/url"
	"whispir/auth-server/services/auth"
)

var (
	RootCmd = cobra.Command{
		Use: "auth",
	}
	RootOpts struct {
		Scheme string
		osincli.ClientConfig
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&RootOpts.Scheme,
		"scheme", "s", "http", "The scheme of the OAuth2 server")
	RootCmd.PersistentFlags().StringVar(&RootOpts.ClientId,
		"client-id", "", "Client id")
	RootCmd.PersistentFlags().StringVar(&RootOpts.ClientSecret,
		"client-secret", "", "Client secret")

}

func preRunE() error {
	if err := checkAuthParams(); nil != err {
		return err
	}
	u := url.URL{
		Scheme: RootOpts.Scheme,
		Host:   "localhost:18081",
		Path:   auth.TokenPath,
	}
	RootOpts.TokenUrl = u.String()

	u.Path = auth.AuthPath
	RootOpts.AuthorizeUrl = u.String()

	RootOpts.RedirectUrl = RedirectURL
	return nil
}

func checkAuthParams() error {
	switch {
	case "" == RootOpts.ClientId:
		return fmt.Errorf("client-id should be specified")
	case "" == RootOpts.ClientSecret:
		return fmt.Errorf("client-sercret should be specified")
	}
	return nil
}
