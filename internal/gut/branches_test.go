package gut

import (
	"reflect"
	"testing"
)

func Test_parseBranchesOutput01(t *testing.T) {
	rawInput :=
		`  XYZ-2252
`
	parsedBranches := parseBranchesOutput([]byte(rawInput))

	tobe := []Branch{{
		Name:      "XYZ-2252",
		Remote:    "",
		IsCurrent: false,
		PointsTo:  "",
	}}

	if !reflect.DeepEqual(tobe, parsedBranches) {
		t.Errorf("tobe: %v + is: %v", tobe, parsedBranches)
	}
}

func Test_parseBranchesOutput02(t *testing.T) {
	rawInput :=
		`  XYZ-2252
* develop
  remotes/origin/HEAD -> origin/master
  remotes/origin/OrderAssistantCleanup
`
	parsedBranches := parseBranchesOutput([]byte(rawInput))

	tobe := []Branch{{
		Name:      "XYZ-2252",
		Remote:    "",
		IsCurrent: false,
		PointsTo:  "",
	}, {
		Name:      "develop",
		Remote:    "",
		IsCurrent: true,
		PointsTo:  "",
	}, {
		Name:      "HEAD",
		Remote:    "origin",
		IsCurrent: false,
		PointsTo:  "origin/master",
	}, {
		Name:      "OrderAssistantCleanup",
		Remote:    "origin",
		IsCurrent: false,
		PointsTo:  "",
	}}

	if !reflect.DeepEqual(tobe, parsedBranches) {
		t.Errorf("parseOutput() tobe: %v + is: %v", tobe, parsedBranches)
	}
}

