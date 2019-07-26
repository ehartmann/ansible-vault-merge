package ansible

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	vault "github.com/sosedoff/ansible-vault-go"
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

func DecryptFile(path string, password string) {
	content, err := vault.DecryptFile(path, password)
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", path, err.Error())
	}
	ioutil.WriteFile(path, []byte(content), 0644)
}

func EncryptFile(path string, password string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", path, err.Error())
	}
	vault.EncryptFile(path, string(content), password)
}

func ViewFile(path string, password string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", path, err.Error())
	}

	decryptedContent, err := vault.Decrypt(string(content), password)
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", path, err.Error())
	}
	return decryptedContent
}

func FindPassword() string {
	ansibleCfgFile := findAnsibleCfg()
	cfg, err := ini.Load(ansibleCfgFile)
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", ansibleCfgFile, err.Error())
	}

	passwordFile := cfg.Section(defaults).Key(vaultPasswordFile).String()
	if len(passwordFile) == 0 {
		log.Fatalf("Did not find vault_password_file inside [%s]", ansibleCfgFile)
	}
	log.Printf("Found passwordFile %s", passwordFile)

	content, err := ioutil.ReadFile(os.ExpandEnv(passwordFile))
	if err != nil {
		log.Fatalf("Unable to read file [%s] : %s", passwordFile, err.Error())
	}

	return string(content[:len(content)-1])
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
