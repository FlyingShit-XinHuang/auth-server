package main

import (
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	httptransport "whispir/auth-server/pkg/transport/http"
	"whispir/auth-server/services/auth"
	"whispir/auth-server/services/client"
	"whispir/auth-server/services/user"
	"whispir/auth-server/storage/mysql"
)

var (
	RootCmd = cobra.Command{
		Use: "auth-server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer()
		},
	}
	RootOpts struct {
		dbhost   string
		password string
		user     string
		port     int
		dbname   string
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&RootOpts.dbhost,
		"host", "H", "localhost", "Database host")
	RootCmd.PersistentFlags().IntVarP(&RootOpts.port,
		"port", "P", 3306, "Database port")
	RootCmd.PersistentFlags().StringVarP(&RootOpts.user,
		"user", "u", "root", "Database user")
	RootCmd.PersistentFlags().StringVarP(&RootOpts.password,
		"password", "p", "", "Database password")
	RootCmd.PersistentFlags().StringVarP(&RootOpts.dbname,
		"db", "d", "demo", "Database name")
}

func main() {
	if err := RootCmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func startServer() error {
	store, err := mysql.NewStorage(RootOpts.user, RootOpts.password, RootOpts.dbhost, RootOpts.port, RootOpts.dbname)
	if nil != err {
		return err
	}

	authSvc := auth.NewBasicService(store)
	clientSvc := client.NewBasicService(store)
	userSvc := user.NewBasiceService(store)

	r := httptransport.NewRouter()
	r.AddHandlers(auth.NewHTTPTransport(authSvc))
	r.AddHandlers(client.NewHTTPTransport(clientSvc))
	r.AddHandlers(user.NewHTTPTransport(userSvc))

	log.Println("Listening 18080")
	log.Fatal(http.ListenAndServe(":18080", r))

	return nil

}