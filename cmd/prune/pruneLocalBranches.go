package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codecare/gut/internal/gut"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("plaese add more arguments: Project path")
		os.Exit(1)
	}
	var baseDir = os.Args[1]
	targetBranch, err := gut.GetDefaultBranch(baseDir)
	if err != nil {
		panic(err)
	}

	projectUrl, err := gut.GetProjectURL(baseDir)
	if err != nil {
		panic(err)
	}

	// fetch and pull changes to update project
	gut.FetchAndPull(baseDir)

	// get all branches (local + remote) & print branch view
	branches, err := gut.ReadBranches(baseDir)
	if err != nil {
		panic(err)
	}
	fmt.Print("List of local branches of the project with url: " + projectUrl + "\n")
	gut.PrintBranches(branches)

	// analyse Merging commits to decide which local branch to delete
	branchesToDelete, err := analyseMerged(branches, baseDir, targetBranch)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\n\nThe following branches are merged to %s and therefore can be deleted: \n", targetBranch)
		for i, branch := range branchesToDelete {
			fmt.Println(strconv.Itoa(i+1) + "- " + branch)
		}
		fmt.Print("\nType numbers of branches from the list you want to delete or type 'a' or 'all' to delete all branches: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if input.Text() == "a" || input.Text() == "all" {
			// delete all local branches
			gut.DeleteLocalBranches(baseDir, branchesToDelete)
		} else {
			var array []string
			for _, value := range strings.Fields(input.Text()) {
				if index, err := strconv.Atoi(value); err == nil && index <= len(branchesToDelete) {
					array = append(array, branchesToDelete[index-1])
				} else {
					fmt.Print("Input error, please type numbers from the list above\n")
					// reset array in case of input error so nothing will be deleted
					array = []string{}
					break
				}
			}
			// delete selected local branches
			gut.DeleteLocalBranches(baseDir, array)
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
