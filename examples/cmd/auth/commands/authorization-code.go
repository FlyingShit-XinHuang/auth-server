package commands

import (
	"fmt"
	"github.com/RangelReale/osincli"
	"github.com/spf13/cobra"
	"net/http"
)

var (
	authCodeCmd = cobra.Command{
		Use:   "ac [flags]",
		Short: "Get access token with authorization code grant type",
		Long: "Get access token with authorization code grant type. " +
			"This command will start a server to serve a redirect url to receive the authorization code " +
			"which could be exchanged for the access token",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := preRunE(); nil != err {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTokenWithAuthCode()
		},
	}
)

func init() {
	RootCmd.AddCommand(&authCodeCmd)
}

const (
	redirectPath = "/appauth"
	redirectPort = ":14000"
	RedirectURL  = "http://localhost" + redirectPort + redirectPath
)

func getTokenWithAuthCode() error {
	client, err := osincli.NewClient(&RootOpts.ClientConfig)
	if nil != err {
		return fmt.Errorf("Failed to init client: %v", err)
	}

	authReq := client.NewAuthorizeRequest(osincli.CODE)

	http.HandleFunc(redirectPath, func(w http.ResponseWriter, r *http.Request) {
		if authCodeData, err := authReq.HandleRequest(r); nil == err {
			data, err := client.NewAccessRequest(osincli.AUTHORIZATION_CODE, authCodeData).GetToken()
			if err == nil {
				w.Write([]byte("Access token: " + data.AccessToken))
			} else {
				w.Write([]byte(fmt.Sprintf("ERROR: %s\n", err)))
			}
		} else {
			w.Write([]byte(fmt.Sprintf("ERROR: %s\n", err)))
		}
	})

	fmt.Println("Please open", authReq.GetAuthorizeUrl().String(), "in your browser")
	http.ListenAndServe(redirectPort, nil)

	return nil
}
