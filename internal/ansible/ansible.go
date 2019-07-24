package ansible

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const ansibleCfg string = "ansible.cfg"
const vaultPasswordFile = "vault_password_file"
const defaults = "defaults"

func visit(ansibleCfgFiles *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == ansibleCfg {
			*ansibleCfgFiles = append(*ansibleCfgFiles, path)
		}

		return nil
	}
}

func execAnsibleVault(passwordFile string, path string, command string) string {
	cmd := exec.Command("sh", "-c", "ansible-vault "+command+" --vault-password-file "+passwordFile+" "+path)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(string(out))
	}

	return string(out)
}

func DecryptFile(passwordFile string, path string) {
	execAnsibleVault(passwordFile, path, "decrypt")
}

func EncryptFile(passwordFile string, path string) {
	execAnsibleVault(passwordFile, path, "encrypt")
}

func ViewFile(passwordFile string, path string) string {
	return execAnsibleVault(passwordFile, path, "view")
}

func FindPasswordFile() string {
	ansibleCfgFile := findAnsibleCfg()
	cfg, err := ini.Load(ansibleCfgFile)
	if err != nil {
		log.Fatal("Fail to read file: %v", err)
	}

	passwordFile := cfg.Section(defaults).Key(vaultPasswordFile).String()
	if len(passwordFile) == 0 {
		log.Fatalf("Did not find vault_password_file inside %s", ansibleCfgFile)
	}
	log.Printf("Found passwordFile %s", passwordFile)

	return passwordFile
}

func findAnsibleCfg() string {
	var ansibleCfg []string
	err := filepath.Walk(".", visit(&ansibleCfg))

	if err != nil {
		panic(err)
	}

	if len(ansibleCfg) > 1 {
		log.Fatal("Found multiple ansible.cfg files")
	} else if len(ansibleCfg) == 0 {
		log.Fatal("Did not found ansible.cfg")
	}

	log.Printf("Found ansible.cfg %s", ansibleCfg[0])

	return ansibleCfg[0]
}
