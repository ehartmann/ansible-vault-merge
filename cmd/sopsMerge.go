package cmd

import (
	"github.com/spf13/cobra"
)

var sopsMergeCmd = &cobra.Command{
	Use:   "sops-merge",
	Short: "Used as a git merge driver",
	Args:  cobra.ExactArgs(3),
	Run:   sopsMerge,
}

func init() {
	RootCmd.AddCommand(sopsMergeCmd)
}

func sopsMerge(cmd *cobra.Command, args []string) {

}
