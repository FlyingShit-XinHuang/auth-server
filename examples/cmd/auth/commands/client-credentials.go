package commands

import (
	"fmt"
	"github.com/RangelReale/osincli"
	"github.com/spf13/cobra"
)

var (
	clientCredentialsCmd = cobra.Command{
		Use:          "cc [flags]",
		Short:        "Get access token with Client Credentials grant type",
		Long:         "Get access token with Client Credentials grant type",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := preRunE(); nil != err {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTokenWithClientCredentials()
		},
	}
)

func init() {
	RootCmd.AddCommand(&clientCredentialsCmd)
}

func getTokenWithClientCredentials() error {
	client, err := osincli.NewClient(&RootOpts.ClientConfig)
	if nil != err {
		return fmt.Errorf("Failed to init client: %v", err)
	}

	data, err := client.NewAccessRequest(osincli.CLIENT_CREDENTIALS, nil).GetToken()

	if nil != err {
		return fmt.Errorf("Failed to get access token: %v\n", err)
	}

	fmt.Println("Access token:", data.AccessToken)

	return nil
}
