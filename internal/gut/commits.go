package gut

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Commit struct {
	Sha1          string
	Parent        []string
	Author        string
	Committer     string
	CommitMessage []string
}

func (commit Commit) String() string {
	return fmt.Sprintf("%s: parents %v, author: %s, committer:%s, msg: %v", commit.Sha1, commit.Parent, commit.Author, commit.Committer, commit.CommitMessage)
}

func GetAllCommits(baseDir, branchName string) ([]Commit, error) {

	command := exec.Command("git", "log", "--pretty=raw", branchName)

	command.Dir = baseDir
	command.Env = append(os.Environ(), "LANG=en_US.UTF-8") // go likes utf-8

	out, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	} else {
		commits, err := parseRawCommits(out)
		return commits, err
	}
}

func GetTopCommit(baseDir, branchName string) (Commit, error) {

	command := exec.Command("git", "log", "-1", "--pretty=raw", branchName)

	command.Dir = baseDir
	command.Env = append(os.Environ(), "LANG=en_US.UTF-8") // go likes utf-8

	out, err := command.CombinedOutput()
	if err != nil {
		return Commit{}, err
	} else {
		commits, err := parseRawCommits(out)
		if err != nil {
			return Commit{}, err
		}
		if len(commits) == 1 {
			return commits[0], nil
		}
		return Commit{}, errors.New("multiple commits found")
	}
}

func PrintCommits(commits []Commit) {
	for _, element := range commits {
		fmt.Printf("commit: %s\n", element)
	}
}

// parses raw commit log (i.e. output of "git log --format=raw") and returns
// an array of Commit objects
//
// sample commit log:
// commit 73700d773cf248834510f6cd06ca7544c934e5df
// tree 042106fec83877a05c3955152e286ea1189f1f6c
// parent d367cbf1ec98bed63345ec000c810e80bba0ad64
// parent 29f9e9f6e4bdb166cfebb417aabc5af5d1b5869a
// author Alex Traud <at@codecare.de> 1568816665 +0200
// committer Alex Traud <at@codecare.de> 1568816665 +0200
//
//    Merge branch 'FEAT1'

func parseRawCommits(bytes []byte) ([]Commit, error) {

	stringData := string(bytes[:])
	scanner := bufio.NewScanner(strings.NewReader(stringData))
	var commits []Commit

	var commit *Commit
	commit = &Commit{}	// just to be sure that we do not run into nil

	for scanner.Scan() {
		s := scanner.Text()

		if strings.Index(s, "fatal:") == 0 {
			return commits, errors.New("s")
		}

		// new commit always starts with commit
		if strings.Index(s, "commit ") == 0 {
			commits = append(commits, Commit{})
			commit = &commits[len(commits) - 1]

			commit.Sha1 = s[len("commit "):]

		} else if strings.Index(s, "parent") == 0 {
			parent := s[len("parent "):]
			commit.Parent = append(commit.Parent, parent)

		} else if strings.Index(s, "author") == 0 {
			commit.Author = s[len("author "):]

		} else if strings.Index(s, "committer ") == 0 {
			commit.Committer = s[len("committer "):]

		} else if strings.Index(s, "    ") == 0 {
			s = strings.TrimLeft(s, " ")
			commit.CommitMessage = append(commit.CommitMessage, s)
		}
	}
	return commits, nil
}

