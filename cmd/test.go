package cmd

import (
	"github.com/spf13/cobra"
)

var consumerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Consumer Test Service CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
