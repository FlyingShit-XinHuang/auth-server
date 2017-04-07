package commands

import (
	"github.com/spf13/cobra"

	"fmt"

	"net/url"
	"whispir/auth-server/examples/cmd/auth/commands"
	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/services/client"
	"whispir/auth-server/services/user"
)

var (
	createCmd = cobra.Command{
		Use:          "create <resource> [flags]",
		Short:        "Create a resource",
		Long:         "Create specified resource",
		SilenceUsage: true,
	}

	clientCmd = cobra.Command{
		Use:          "client [flags] <name>",
		SilenceUsage: true,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing arguments")
			}
			clientOpts.client = &v1alpha1.Client{
				Name:        args[0],
				RedirectURL: commands.RedirectURL,
			}
			serviceURI.Path = client.CreateClientPath
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			newclient, err := client.NewServiceForHTTPClient(serviceURI).CreateClient(clientOpts.client)
			if nil != err {
				return fmt.Errorf("Failed to create client: %v", err)
			}
			fmt.Printf("Created client with id '%s' and secret '%s'\n", newclient.Id, newclient.Secret)
			return nil
		},
	}

	userCmd = cobra.Command{
		Use:          "user [flags] <name> <password>",
		SilenceUsage: true,
		PreRunE: func(_ *cobra.Command, args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("missing arguments")
			}
			userOpts.user = &v1alpha1.User{
				Name:     args[0],
				Password: args[1],
			}
			serviceURI.Path = user.CreateUserPath
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			err := user.NewServiceForHTTPClient(serviceURI).CreateUser(userOpts.user)
			if nil != err {
				return fmt.Errorf("Failed to create user: %v", err)
			}
			fmt.Println("Success")
			return nil
		},
	}

	clientOpts struct {
		client *v1alpha1.Client
	}

	userOpts struct {
		user *v1alpha1.User
	}

	serviceURI *url.URL
)

func init() {
	var err error
	serviceURI, err = url.Parse("http://localhost:18080")
	if nil != err {
		panic(err)
	}

	RootCmd.AddCommand(&createCmd)

	createCmd.AddCommand(&clientCmd)
	createCmd.AddCommand(&userCmd)
}
