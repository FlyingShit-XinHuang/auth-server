package commands

import (
	//"flag"
	"fmt"
	"github.com/RangelReale/osincli"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"whispir/auth-server/services/auth"
)

var (
	RootCmd = cobra.Command{
		Use: "auth",
	}
	RootOpts struct {
		Scheme string
		osincli.ClientConfig
		host       string
		port       int
		pathPrefix string
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&RootOpts.Scheme,
		"scheme", "s", "http", "The scheme of the OAuth2 server")
	RootCmd.PersistentFlags().StringVar(&RootOpts.ClientId,
		"client-id", "", "Client id")
	RootCmd.PersistentFlags().StringVar(&RootOpts.ClientSecret,
		"client-secret", "", "Client secret")

	RootCmd.PersistentFlags().StringVarP(&RootOpts.host,
		"host", "H", "localhost", "auth-server host")
	RootCmd.PersistentFlags().IntVarP(&RootOpts.port,
		"port", "P", 18080, "auth-server port")
	RootCmd.PersistentFlags().StringVarP(&RootOpts.pathPrefix,
		"prefix", "p", "/", "Prefix of path of request")

}

func preRunE() error {
	if err := checkAuthParams(); nil != err {
		return err
	}
	u := url.URL{
		Scheme: RootOpts.Scheme,
		Host:   fmt.Sprintf("%s:%d", RootOpts.host, RootOpts.port),
		Path:   path.Join(RootOpts.pathPrefix, auth.TokenPath),
	}
	RootOpts.TokenUrl = u.String()

	u.Path = path.Join(RootOpts.pathPrefix, auth.AuthPath)
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
