package main

import (
	"log"
	"net/http"
	httptransport "whispir/auth-server/pkg/transport/http"
	"whispir/auth-server/services/auth"
	"whispir/auth-server/services/client"
	"whispir/auth-server/services/user"
	"whispir/auth-server/storage/mysql"
)

func main() {
	store := mysql.NewStorageOrDie("root", "demo", "localhost", 3306, "demo")
	authSvc := auth.NewBasicService(store)
	clientSvc := client.NewBasicService(store)
	userSvc := user.NewBasiceService(store)


	r := httptransport.NewRouter()
	r.AddHandlers(auth.NewHTTPTransport(authSvc))
	r.AddHandlers(client.NewHTTPTransport(clientSvc))
	r.AddHandlers(user.NewHTTPTransport(userSvc))

	log.Println("Listening 18080")
	log.Fatal(http.ListenAndServe(":18080", r))
}
