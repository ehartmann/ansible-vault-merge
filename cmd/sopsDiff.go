package cmd

import (
	"github.com/spf13/cobra"
)

var sopsDiffCmd = &cobra.Command{
	Use:   "sops-diff",
	Short: "Used as a git diff driver",
	Args:  cobra.ExactArgs(3),
	Run:   sopsDiff,
}

func init() {
	RootCmd.AddCommand(sopsDiffCmd)
}

func sopsDiff(cmd *cobra.Command, args []string) {

}
