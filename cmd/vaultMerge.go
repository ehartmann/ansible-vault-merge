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
var editor string
var stop bool

func init() {
	vaultMergeCmd.Flags().BoolVarP(&stop, "stop", "s", false, "Stop if there is a conflict")
	vaultMergeCmd.Flags().StringVarP(&editor, "editor", "e", "", "Editor to use for git conflict")
	RootCmd.AddCommand(vaultMergeCmd)
}

func vaultMerge(cmd *cobra.Command, args []string) {
	base := args[0]
	current := args[1]
	other := args[2]

	passwordFile := ansible.FindPasswordFile()

	ansible.DecryptFile(passwordFile, base)
	ansible.DecryptFile(passwordFile, current)
	ansible.DecryptFile(passwordFile, other)

	err := git.GitMerge(base, current, other)
	log.Print(err.Error())
	if err != nil {
		if stop {
			log.Fatal("Conflicts detected. Stopping, do not forgot to encrypt file !")
		}
		utils.LaunchEditor(current, editor)
	}

	if !git.IsAllConflictsSolved(current) {
		log.Fatal("Conflict not solved. Stopping, do not forgot to encrypt file !")
	}

	ansible.EncryptFile(passwordFile, current)
}
