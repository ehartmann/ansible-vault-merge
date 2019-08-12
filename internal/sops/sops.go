package sops

import (
	"io/ioutil"
	"log"
	"os/exec"

	"go.mozilla.org/sops/decrypt"
)

func Decrypt(path string) (string, error) {
	// FIXME: during merge temporary files are named like .merge_file_xxx
	// so we cannot determine the file type by using the suffix, however
	// we are always using yaml for encrypted files
	data, err := decrypt.File(path, "yaml")

	return string(data), err
}

func DecryptFile(path string) error {
	data, err := Decrypt(path)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(data), 0644)
}

func EncryptFile(path string) error {
	cmd := exec.Command("sh", "-c", "sops -e "+path)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(out))
	}
	return err
}
