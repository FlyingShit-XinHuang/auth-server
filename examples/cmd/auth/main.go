package main

import (
	"os"
	"whispir/auth-server/examples/cmd/auth/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); nil != err {
		os.Exit(1)
	}
}
