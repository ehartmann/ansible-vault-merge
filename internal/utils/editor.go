package utils

import (
	"log"
	"errors"
	"os"
	"os/exec"
)

func LaunchEditor(path string, editor string) error {
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		return errors.New("No editor configured")
	}
	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdoutgit s
	editorCmd.Stderr = os.Stderr
	editorErr := editorCmd.Start()
	if editorErr != nil {
		return editorErr
	}
	err := editorCmd.Wait()
	return err
}
