package testin

import "testing"

func TestSplitParts(t *testing.T) {
	parts, err := splitParts(`main
########
# FOO
########
FOO
########
# BAR
########
BAR`)
	if err != nil {
		t.Fatal(err)
	}
	if parts.main != "main" {
		t.Errorf("Main part should be main but was '%v'", parts.main)
	}
	if parts.results["FOO"] != "FOO" {
		t.Errorf("FOO part should be 'foo' but was '%v'", parts.results["FOO"])
	}
	if parts.results["BAR"] != "BAR" {
		t.Errorf("BAR part should be 'bar' but was '%v'", parts.results["BAR"])
	}
}
