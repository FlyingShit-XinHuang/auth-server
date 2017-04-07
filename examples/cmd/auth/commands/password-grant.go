package commands

import (
	"fmt"
	"github.com/RangelReale/osincli"
	"github.com/spf13/cobra"
)

var (
	passwordGrantCmd = cobra.Command{
		Use:          "pg [flags] <username> <password>",
		Short:        "Get access token with Password grant type",
		Long:         "Get access token with Password grant type",
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := preRunE(); nil != err {
				return err
			}
			if len(args) < 2 {
				return fmt.Errorf("Missing args")
			}
			passwordGrantOpts.user = args[0]
			passwordGrantOpts.password = args[1]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTokenWithPasswordGrant()
		},
	}

	passwordGrantOpts struct {
		user     string
		password string
	}
)

func init() {
	RootCmd.AddCommand(&passwordGrantCmd)
}

func getTokenWithPasswordGrant() error {
	if "" == passwordGrantOpts.user {
		return fmt.Errorf("username should be specified")
	}
	if "" == passwordGrantOpts.password {
		return fmt.Errorf("password should be specified")
	}

	client, err := osincli.NewClient(&RootOpts.ClientConfig)
	if nil != err {
		return fmt.Errorf("Failed to init client: %v", err)
	}

	data, err := client.NewAccessRequest(osincli.PASSWORD, &osincli.AuthorizeData{
		Username: passwordGrantOpts.user,
		Password: passwordGrantOpts.password,
	}).GetToken()

	if nil != err {
		return fmt.Errorf("Failed to get access token: %v\n", err)
	}

	fmt.Println("Access token:", data.AccessToken)

	return nil
}
