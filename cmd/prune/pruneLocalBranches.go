package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/codecare/gut/internal/gut"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("plaese add more arguments: Project path + branch name")
		os.Exit(1)
	}
	var baseDir = os.Args[1]
	var targetBranch = os.Args[2]

	// fetch and pull changes to update project
	gut.FetchAndPull(baseDir)

	// get all branches (local + remote) & print branch view
	branches, err := gut.ReadBranches(baseDir)
	if err != nil {
		panic(err)
	}
	gut.PrintBranches(branches)

	// analyse Merging commits to decide which local branch to delete
	branchesToDelete, err := analyseMerged(branches, baseDir, targetBranch)
	if err != nil {
		panic(err)
	} else {
		var output = strings.Join(branchesToDelete, `, `)
		fmt.Printf("The following branches are merged to %s and therefore can be deleted \n%s\n", targetBranch, output)
		fmt.Print("Do you really want to delete them? Type y or yes to confirm: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if input.Text() == "y" || input.Text() == "yes" {
			// delete local branches
			gut.DeleteLocalBranches(baseDir, branchesToDelete)
		}
	}
}

func analyseMerged(branches []gut.Branch, baseDir, targetBranchName string) ([]string, error) {

	allCommitsOfTarget, err := gut.GetAllCommits(baseDir, targetBranchName)
	if err != nil {
		panic(err)
	}
	var branchesToDelete []string

	for _, branch := range branches {

		// skip current branch
		if branch.IsCurrent {
			continue
		}
		// skip remote branches, they will be deleted separately
		if branch.Remote != "" {
			continue
		}
		branchMergeAnalysis, err := analyseBranch(branch, targetBranchName, baseDir, allCommitsOfTarget)
		if err != nil {
			return branchesToDelete, err
		}
		if branchMergeAnalysis == FullyMerged {
			branchesToDelete = append(branchesToDelete, branch.Name)
		}
	}
	return branchesToDelete, nil
}

func analyseBranch(branchToAnalyse gut.Branch, targetBranchName string, baseDir string, allCommitsOfTarget []gut.Commit) (MergeType, error) {

	if branchToAnalyse.Name == targetBranchName {
	} else {
		topCommitOfBranch, err := gut.GetTopCommit(baseDir, branchToAnalyse.Name)
		if err != nil {
			return NotMerged, err
		}
		err = findMergeCommit(allCommitsOfTarget, topCommitOfBranch)
		if err != nil {
			return PartiallyMerged, nil
		} else {
			return FullyMerged, nil
		}
	}
	return SameBranch, nil
}

func findMergeCommit(allCommitsOfTarget []gut.Commit, topCommitOfBranch gut.Commit) error {
	for _, commit := range allCommitsOfTarget {
		if isParentCommit(topCommitOfBranch, commit) {
			return nil
		}
	}
	return errors.New("merge commit not found")
}

func isParentCommit(branchToCheck gut.Commit, potentialParent gut.Commit) bool {
	for _, sha1 := range potentialParent.Parent {
		if sha1 == branchToCheck.Sha1 {
			return true
		}
	}
	return false
}

type MergeType int

const (
	FullyMerged MergeType = iota
	PartiallyMerged
	SameBranch
	NotMerged
)
