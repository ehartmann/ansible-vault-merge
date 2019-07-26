package git

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

var conflictMarkers = []string{
	"<<<<<<< ",
	"======= ",
	"=======\n",
	">>>>>>> "}

func IsAllConflictsSolved(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		for _, marker := range conflictMarkers {
			if strings.HasPrefix(scanner.Text(), marker) {
				return false
			}
		}
	}

	return true
}

func GitMerge(base string, current string, other string) error {
	cmd := exec.Command("sh", "-c", "git merge-file -L CURRENT -L BASE -L OTHER "+current+" "+base+" "+other)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(out))
	}
	return err
}

func GitAdd(path string) error {
	cmd := exec.Command("sh", "-c", "git add "+path)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Print(string(out))
	}
	return err
}
