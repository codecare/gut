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
		index = strings.Index(branchName, "remotes/")
		if index == 0 {
			branchName = branchName[8:]
			index = strings.Index(branchName, "/")
			if index > 0 {
				remote = branchName[:index]
				branchName = branchName[index+1:]
			} else {
				panic("remote cannot be parsed " + branchName)
			}
		}

		branch := Branch{Name: branchName, IsCurrent: isCurrent, PointsTo: pointsTo, Remote: remote}
		branches = append(branches, branch)
	}
	return branches
}
