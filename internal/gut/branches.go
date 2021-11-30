package gut

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Branch struct {
	Name      string
	Remote    string
	IsCurrent bool
	PointsTo  string
}

func (branch Branch) String() string {
	return fmt.Sprintf("%s, isCurrent: %t, pointsTo: %s, remote:%s", branch.Name, branch.IsCurrent, branch.PointsTo, branch.Remote)
}

func FetchAndPull(baseDir string) {

	command := exec.Command("git", "fetch", "-p")
	command.Dir = baseDir
	command.Env = append(os.Environ(),
		"LANG=en_US.UTF-8")

	out, err := command.CombinedOutput()
	if err != nil {
		fmt.Printf("Couldn't fetch, error: %s\n", err)
	} else {
		fmt.Printf("successfully fetch command executed: %s\n", out)
	}

	command2 := exec.Command("git", "pull")
	command2.Dir = baseDir
	command2.Env = append(os.Environ(),
		"LANG=en_US.UTF-8")

	out, err = command2.CombinedOutput()
	if err != nil {
		fmt.Printf("Couldn't pull, error: %s\n", err)
	} else {
		fmt.Printf("successfully pull command executed: %s\n", out)
	}
}

func DeleteLocalBranches(baseDir string, branchesToDelete []string) {

	for _, branch := range branchesToDelete {
		command := exec.Command("git", "branch", "-d", branch)
		command.Dir = baseDir
		command.Env = append(os.Environ(),
			"LANG=en_US.UTF-8") // go likes utf-8

		out, err := command.CombinedOutput()
		if err != nil {
			fmt.Printf("Couldn't delete branch %s, error: %s\n", branch, err)
		} else {
			fmt.Printf("Deleting local branch %s: %s\n", branch, out)
		}
	}
}

func ReadBranches(baseDir string) ([]Branch, error) {

	command := exec.Command("git", "branch", "-a")
	// todo uli should be current dir
	command.Dir = baseDir
	command.Env = append(os.Environ(),
		"LANG=en_US.UTF-8") // go likes utf-8

	out, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	} else {
		branches := parseBranchesOutput(out)
		return branches, nil
	}
}

func PrintBranches(Branches []Branch) {
	for _, element := range Branches {
		fmt.Printf("branch: %s\n", element)
	}
}

func parseBranchesOutput(bytes []byte) []Branch {
	// go is automatically utf-8
	stringData := string(bytes[:])
	scanner := bufio.NewScanner(strings.NewReader(stringData))
	var branches []Branch
	for scanner.Scan() {
		// fmt.Printf("line: %s\n", scanner.Text())
		branchName := strings.TrimSpace(scanner.Text())

		isCurrent := false
		if strings.Index(branchName, "*") == 0 {
			isCurrent = true
			branchName = strings.TrimSpace(branchName[1:])
		}

		pointsTo := ""
		index := strings.Index(branchName, "->")
		if index > 0 {
			pointsTo = strings.TrimSpace(branchName[index+2:])
			branchName = strings.TrimSpace(branchName[:index])
		}

		remote := ""
		var fullBranchName = branchName
		index = strings.Index(branchName, "remotes/")
		if index == 0 {
			branchName = branchName[8:]
			index = strings.Index(branchName, "/")
			if index > 0 {
				remote = branchName[:index]
				branchName = fullBranchName
			} else {
				panic("remote cannot be parsed " + branchName)
			}
		}

		branch := Branch{Name: branchName, IsCurrent: isCurrent, PointsTo: pointsTo, Remote: remote}
		branches = append(branches, branch)
	}
	return branches
}
