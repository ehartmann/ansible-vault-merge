package cmd

import (
	"fmt"

	"github.com/ehartmann/merge-secrets/internal/ansible"

	"github.com/spf13/cobra"
)

var vaultDiffCmd = &cobra.Command{
	Use:   "vault-textconv",
	Short: "Used as a git diff textconv",
	Args:  cobra.ExactArgs(1),
	Run:   vaultDiff,
}

func init() {
	RootCmd.AddCommand(vaultDiffCmd)
}

func vaultDiff(cmd *cobra.Command, args []string) {
	path := args[0]

	password := ansible.FindPassword()

	fmt.Print(ansible.ViewFile(path, password))
}
