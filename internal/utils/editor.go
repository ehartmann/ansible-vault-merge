package utils

import (
	"log"
	"os"
	"os/exec"
)

func LaunchEditor(path string, editor string) {
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	editorErr := editorCmd.Start()
	if editorErr != nil {
		log.Fatal(editorErr.Error())
	}
	editorCmd.Wait()
}
