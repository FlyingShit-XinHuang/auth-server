package main

import (
	"os"
	"whispir/auth-server/examples/cmd/resources/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); nil != err {
		os.Exit(1)
	}
}
