package cmd

import (
	"log"

	"github.com/ehartmann/merge-secrets/internal/git"
	"github.com/ehartmann/merge-secrets/internal/sops"
	"github.com/ehartmann/merge-secrets/internal/utils"
	"github.com/spf13/cobra"
)

var sopsMergeCmd = &cobra.Command{
	Use:   "sops-merge",
	Short: "Used as a git merge driver",
	Args:  cobra.ExactArgs(3),
	Run:   sopsMerge,
}

func init() {
	sopsMergeCmd.Flags().BoolVarP(&stop, "stop", "s", false, "Stop if there is a conflict")
	sopsMergeCmd.Flags().StringVarP(&editor, "editor", "e", "", "Editor to use for git conflict")
	RootCmd.AddCommand(sopsMergeCmd)
}

func sopsMerge(cmd *cobra.Command, args []string) {
	base := args[0]
	current := args[1]
	other := args[2]

	err := sops.DecryptFile(base)
	if err != nil {
		log.Fatalf("Unable to decrypt file %s : %s", base, err)
	}
	err = sops.DecryptFile(current)
	if err != nil {
		log.Fatalf("Unable to decrypt file %s : %s", current, err)
	}
	err = sops.DecryptFile(other)
	if err != nil {
		log.Fatalf("Unable to decrypt file %s : %s", other, err)
	}

	err = git.GitMerge(base, current, other)
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

	sops.EncryptFile(current)
	if err != nil {
		log.Fatalf("Unable to encrypt file %s", err)
	}
}
