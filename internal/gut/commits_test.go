package gut

import (
	"reflect"
	"testing"
)

func Test_parseRawCommits01(t *testing.T) {
	rawInput :=
		`commit f96c00eea6501fcfbd5f7cb6e18ef3e6d4a3b7c8
tree 2087748f7f7e6ba07b595d2f2527e392fbc07755
parent a041ca9099d549630ece5705c807dbb530fba14a
parent 7ca7db67e84ebbad40b302e2fc82e89b5b212ff1
author Alex Traud <at@codecare.de> 1567596532 +0200
committer Alex Traud <at@codecare.de> 1567596532 +0200

    Merge branch 'XYZ-2094-braintree-payment-non-euro' into develop

`
	parsedCommits, err := parseRawCommits([]byte(rawInput))
	if !(err == nil) {
		t.Errorf("error: %v\n", err)
	}

	tobe := []Commit{{
		Sha1:          "f96c00eea6501fcfbd5f7cb6e18ef3e6d4a3b7c8",
		Parent:        []string{"a041ca9099d549630ece5705c807dbb530fba14a", "7ca7db67e84ebbad40b302e2fc82e89b5b212ff1"},
		Author:        "Alex Traud <at@codecare.de> 1567596532 +0200",
		Committer:     "Alex Traud <at@codecare.de> 1567596532 +0200",
		CommitMessage: []string{"Merge branch 'XYZ-2094-braintree-payment-non-euro' into develop"},
	}}

	if !reflect.DeepEqual(tobe, parsedCommits) {
		t.Errorf("tobe: %v + \n is: %v", tobe, parsedCommits)
	}
}

func Test_parseRawCommits02(t *testing.T) {
	rawInput :=
		`commit f96c00eea6501fcfbd5f7cb6e18ef3e6d4a3b7c8
tree 2087748f7f7e6ba07b595d2f2527e392fbc07755
parent a041ca9099d549630ece5705c807dbb530fba14a
parent 7ca7db67e84ebbad40b302e2fc82e89b5b212ff1
author Alex Traud <at@codecare.de> 1567596532 +0200
committer Alex Traud <at@codecare.de> 1567596532 +0200

    Merge branch 'XYZ-2094-braintree-payment-non-euro' into develop

commit 7ca7db67e84ebbad40b302e2fc82e89b5b212ff1
tree 2c0b4e559ecaae52be8eea4a797c8e70c76b0d8f
parent e73af741547f22b2c5696b07630db31c9341b895
author Alex Traud <at@codecare.de> 1567596263 +0200
committer Alex Traud <at@codecare.de> 1567596263 +0200

    XYZ-2094 TMO | Upgrade to 3D Secure 2.0 authentication protocol (Braintree)
    
    fixed missing conversion to EUR for other currencies for braintree
    based payments.

`
	parsedCommits, err := parseRawCommits([]byte(rawInput))
	if !(err == nil) {
		t.Errorf("error: %v\n", err)
	}

	tobe := []Commit{{
		Sha1:          "f96c00eea6501fcfbd5f7cb6e18ef3e6d4a3b7c8",
		Parent:        []string{"a041ca9099d549630ece5705c807dbb530fba14a", "7ca7db67e84ebbad40b302e2fc82e89b5b212ff1"},
		Author:        "Alex Traud <at@codecare.de> 1567596532 +0200",
		Committer:     "Alex Traud <at@codecare.de> 1567596532 +0200",
		CommitMessage: []string{"Merge branch 'XYZ-2094-braintree-payment-non-euro' into develop"},
	}, {
		Sha1:          "7ca7db67e84ebbad40b302e2fc82e89b5b212ff1",
		Parent:        []string{"e73af741547f22b2c5696b07630db31c9341b895"},
		Author:        "Alex Traud <at@codecare.de> 1567596263 +0200",
		Committer:     "Alex Traud <at@codecare.de> 1567596263 +0200",
		CommitMessage: []string{"XYZ-2094 TMO | Upgrade to 3D Secure 2.0 authentication protocol (Braintree)", "", "fixed missing conversion to EUR for other currencies for braintree", "based payments."},
	}}

	if !reflect.DeepEqual(tobe, parsedCommits) {
		t.Errorf("\ntobe: %v + \n  is: %v", tobe, parsedCommits)
	}
}