package cmd

import (
	"log"

	"github.com/ehartmann/merge-secrets/internal/ansible"
	"github.com/ehartmann/merge-secrets/internal/git"
	"github.com/ehartmann/merge-secrets/internal/utils"
	"github.com/spf13/cobra"
)

var vaultMergeCmd = &cobra.Command{
	Use:   "vault-merge",
	Short: "Used as a git merge driver",
	Args:  cobra.ExactArgs(3),
	Run:   vaultMerge,
}

func init() {
	vaultMergeCmd.Flags().BoolVarP(&stop, "stop", "s", false, "Stop if there is a conflict")
	vaultMergeCmd.Flags().StringVarP(&editor, "editor", "e", "", "Editor to use for git conflict")
	RootCmd.AddCommand(vaultMergeCmd)
}

func vaultMerge(cmd *cobra.Command, args []string) {
	base := args[0]
	current := args[1]
	other := args[2]

	password := ansible.FindPassword()

	ansible.DecryptFile(base, password)
	ansible.DecryptFile(current, password)
	ansible.DecryptFile(other, password)

	err := git.GitMerge(base, current, other)
	log.Print(err.Error())
	if err != nil {
		if stop {
			log.Fatal("Conflicts detected. Stopping, do not forgot to encrypt file !")
		}

		editorErr := utils.LaunchEditor(current, editor)
		if editorErr != nil {
			log.Fatalf("Error launching editor [%s]. Stopping here, do not forgot to encrypt file", editorErr)
		}

		if !git.IsAllConflictsSolved(current) {
			log.Fatal("Conflict not solved. Stopping, do not forgot to encrypt file !")
		}
	}

	ansible.EncryptFile(current, password)
}
