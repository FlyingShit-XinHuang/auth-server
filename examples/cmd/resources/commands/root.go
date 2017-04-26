package commands

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = cobra.Command{
		Use:   "manage",
		Short: "Manage client and user resource",
	}
)