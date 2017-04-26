package commands

import (
	"github.com/spf13/cobra"

	"fmt"

	"net/url"
	"whispir/auth-server/examples/cmd/auth/commands"
	"whispir/auth-server/pkg/api/v1alpha1"
	"whispir/auth-server/services/client"
	"whispir/auth-server/services/user"
	"path"
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
			if err := parseServerHost(); nil != err {
				return err
			}
			clientOpts.client = &v1alpha1.Client{
				Name:        args[0],
				RedirectURL: commands.RedirectURL,
			}
			serviceURI.Path = path.Join(createOpts.pathPrefix, client.CreateClientPath)
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
			if err := parseServerHost(); nil != err {
				return err
			}
			userOpts.user = &v1alpha1.User{
				Name:     args[0],
				Password: args[1],
			}
			serviceURI.Path = path.Join(createOpts.pathPrefix, user.CreateUserPath)
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

	createOpts struct {
		host string
		port int
		pathPrefix string
	}

	clientOpts struct {
		client *v1alpha1.Client
	}

	userOpts struct {
		user *v1alpha1.User
	}

	serviceURI *url.URL
)

func init()  {
	createCmd.PersistentFlags().StringVarP(&createOpts.host,
		"host", "H", "localhost", "auth-server host")
	createCmd.PersistentFlags().IntVarP(&createOpts.port,
		"port", "P", 18080, "auth-server port")
	createCmd.PersistentFlags().StringVarP(&createOpts.pathPrefix,
		"prefix", "p", "/", "Prefix of path of request")
}

func init() {
	var err error
	if nil != err {
		panic(err)
	}

	RootCmd.AddCommand(&createCmd)

	createCmd.AddCommand(&clientCmd)
	createCmd.AddCommand(&userCmd)
}

func parseServerHost() (err error) {
	serviceURI, err = url.Parse(fmt.Sprintf("http://%s:%d", createOpts.host, createOpts.port))
	return
}