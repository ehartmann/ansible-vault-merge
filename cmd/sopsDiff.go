package cmd

import (
	"fmt"
	"log"

	"github.com/ehartmann/merge-secrets/internal/sops"
	"github.com/spf13/cobra"
)

var sopsDiffCmd = &cobra.Command{
	Use:   "sops-textconv",
	Short: "Used as a git diff driver",
	Args:  cobra.ExactArgs(1),
	Run:   sopsDiff,
}

func init() {
	RootCmd.AddCommand(sopsDiffCmd)
}

func sopsDiff(cmd *cobra.Command, args []string) {
	path := args[0]

	data, err := sops.Decrypt(path)

	if err != nil {
		log.Fatalf("Unable to decrypt SOPS file : %s", err)
	}

	fmt.Print(data)
}
