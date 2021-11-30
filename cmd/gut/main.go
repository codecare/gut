package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/codecare/gut/internal/gut"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("usage: \ngut /path/to/gitdir/ develop\ndevelop is the target branch which should be merged")
		os.Exit(1)
	}
	var baseDir = os.Args[1]
	var targetBranch = os.Args[2]

	branches, err := gut.ReadBranches(baseDir)
	if err != nil {
		panic(err)
	}
	gut.PrintBranches(branches)
	err = analyseMerged(branches, baseDir, targetBranch)

	if err != nil {
		panic(err)
	}
	/*
		commits, err := gut.GetAllCommits(baseDir, "develop")
		if err != nil {
			panic(err)
		}*/
	// gut.PrintCommits(commits)
	// analyseMergeCommits(commits)

}

func analyseMergeCommits(commits []gut.Commit) {
	// TODO: add Branch[] to check if branches are merged and identify merge commit for respective branch
	for _, commit := range commits {
		// fmt.Printf("branch: %s\n", commit)
		if len(commit.Parent) > 1 {
			fmt.Printf("merge commit: %v\n", commit.CommitMessage)
		}
	}
}

func analyseMerged(branches []gut.Branch, baseDir, targetBranchName string) error {

	allCommitsOfTarget, err := gut.GetAllCommits(baseDir, targetBranchName)
	if err != nil {
		panic(err)
	}

	for _, branch := range branches {

		if branch.IsCurrent {
			continue
		}
		branchMergeAnalysis, err := analyseBranch(branch, targetBranchName, baseDir, allCommitsOfTarget)
		if err != nil {
			return err
		}
		fmt.Printf("#########\nanalysis of %s in respect to %s \nanalysis result: %s\n", branch.Name, targetBranchName, branchMergeAnalysis)

	}
	return nil
}

func analyseBranch(branchToAnalyse gut.Branch, targetBranchName string, baseDir string, allCommitsOfTarget []gut.Commit) (BranchMergeAnalysis, error) {

	if branchToAnalyse.Name == targetBranchName {
		fmt.Printf("skipping target branch %s\n", targetBranchName)
	} else {
		topCommitOfBranch, err := gut.GetTopCommit(baseDir, branchToAnalyse.Name)
		if err != nil {
			return BranchMergeAnalysis{}, err
		}
		// fmt.Printf("top commit of branch %s: %s\n", branchToAnalyse.Name, topCommitOfBranch)
		mergeCommit, err := findMergeCommit(allCommitsOfTarget, topCommitOfBranch)
		if err != nil {
			fmt.Printf("### No merge commit found for branch %s\n", branchToAnalyse.Name)
			return BranchMergeAnalysis{
				mergeType:         PartiallyCherryPicked, // actually this is just a guess - TODO
				branchAnalysed:    branchToAnalyse,
				topCommitOfBranch: topCommitOfBranch,
				targetBranchName:  targetBranchName,
			}, nil
		} else {
			fmt.Printf("### %s was merged in commit %s into %s\n", branchToAnalyse.Name, mergeCommit, targetBranchName)
			return BranchMergeAnalysis{
				mergeType:             FullyMerged,
				branchAnalysed:        branchToAnalyse,
				topCommitOfBranch:     topCommitOfBranch,
				targetBranchName:      targetBranchName,
				mergeCommitIntoTarget: mergeCommit,
			}, nil
		}
	}
	return BranchMergeAnalysis{
		mergeType:        SameBranch,
		branchAnalysed:   branchToAnalyse,
		targetBranchName: targetBranchName,
	}, nil
}

func findMergeCommit(allCommitsOfTarget []gut.Commit, topCommitOfBranch gut.Commit) (gut.Commit, error) {
	for _, commit := range allCommitsOfTarget {
		if isParentCommit(topCommitOfBranch, commit) {
			return commit, nil
		}
	}
	return gut.Commit{}, errors.New("merge commit not found")
}

func isParentCommit(branchToCheck gut.Commit, potentialParent gut.Commit) bool {
	for _, sha1 := range potentialParent.Parent {
		if sha1 == branchToCheck.Sha1 {
			return true
		}
	}
	return false
}

type BranchMergeAnalysis struct {
	mergeType             MergeType
	branchAnalysed        gut.Branch
	topCommitOfBranch     gut.Commit
	targetBranchName      string
	mergeCommitIntoTarget gut.Commit
}
type MergeType int

const (
	FullyMerged MergeType = iota
	PartiallyCherryPicked
	SameBranch
	NotMerged
)

func (mergeType MergeType) String() string {
	return [...]string{"FullyMerged", "PartiallyMerged", "SameBranch", "NotMerged"}[mergeType]
}

func (branchMergeAnalysis BranchMergeAnalysis) String() string {

	switch branchMergeAnalysis.mergeType {
	case SameBranch:
		{
			return fmt.Sprintf("%s is same branch as %s\n", branchMergeAnalysis.branchAnalysed.Name, branchMergeAnalysis.targetBranchName)
		}
	case FullyMerged:
		{
			result := fmt.Sprintf("%s %s in %s\n", branchMergeAnalysis.mergeType, branchMergeAnalysis.branchAnalysed.Name, branchMergeAnalysis.targetBranchName)
			result += fmt.Sprintf("merge commit: %v\n", branchMergeAnalysis.mergeCommitIntoTarget)
			result += fmt.Sprintf("top commit of branch: %v\n", branchMergeAnalysis.topCommitOfBranch)
			return result
		}
	case PartiallyCherryPicked:
		{
			return fmt.Sprintf("%s %s in %s\n", branchMergeAnalysis.mergeType, branchMergeAnalysis.branchAnalysed.Name, branchMergeAnalysis.targetBranchName)
		}
	case NotMerged:
		{
			return fmt.Sprintf("%s %s in %s\n", branchMergeAnalysis.mergeType, branchMergeAnalysis.branchAnalysed.Name, branchMergeAnalysis.targetBranchName)
		}

	}
	return ""
}
