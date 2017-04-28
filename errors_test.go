package gutil

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func nilError() error {
	return nil
}

func TestFatalIf(t *testing.T) {
	testCases := []struct {
		name    string
		message string
		params  []interface{}
		crash   bool
	}{
		{
			"error-arg",
			"I should crash because %v",
			[]interface{}{errors.New("invalid")},
			true,
		},
		{
			"non-error-args",
			"I should not crash because only non error parameters, %v, %v, %v",
			[]interface{}{"", 4, "foo"},
			false,
		},
		{
			"non-error-nil-args",
			"I should not crash because nil errors, %v, %v, %v, %v",
			[]interface{}{"", 4, "foo", nilError()},
			false,
		},
		{
			"args-with-err",
			"I should crash because an error, %v, %v, %v, %v, %v",
			[]interface{}{"", 4, fmt.Errorf("errr"), "foo", nilError()},
			true,
		},
	}

	// check if it is executed by itself, if yes, run specified test case
	subCall := os.Getenv("TEST_FATAL_IF")
	if subCall != "" {
		i, _ := strconv.Atoi(subCall)
		testCase := testCases[i]
		FatalIf(testCase.message, testCase.params...)
		return
	}

	for index, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			// run test again with env var set to identify test case
			cmd := exec.Command(os.Args[0], "-test.run=TestFatalIf")
			cmd.Env = append(os.Environ(), fmt.Sprintf("TEST_FATAL_IF=%v", index))
			err := cmd.Run()

			if exit, isExitError := err.(*exec.ExitError); isExitError && !exit.Success() {
				// execution exited erroneous
				if !testCase.crash {
					t.Fatalf("process crashed but should have run successfull")
				}
			} else {
				if testCase.crash {
					t.Fatalf("process ran succesful but should have crashed")
				}
			}
		})
	}
}

func TestErrorCollector(t *testing.T) {
	ec := NewErrorCollector()
	ec.Add(nil)
	ec.Add(nilError())
	ec.Add(errors.New("FIRST"))
	ec.Add(fmt.Errorf("LAST"))

	if !ec.Has() {
		t.Errorf("collector should have two errors")
	}

	if ec.First().Error() != "FIRST" {
		t.Errorf("first error should be FIRST")
	}

	if ec.Last().Error() != "LAST" {
		t.Errorf("last error should be LAST")
	}

	if ec.Error() != "FIRST\nLAST" {
		t.Errorf("errors should be FIRST\\nLAST")
	}
}

func TestEmptyCollector(t *testing.T) {
	ec := NewErrorCollector()
	ec.Add(nil)
	ec.Add(nilError())

	if ec.Has() {
		t.Errorf("collector should have no errors")
	}

	if ec.First() != nil {
		t.Errorf("first error should be nil")
	}

	if ec.Last() != nil {
		t.Errorf("last error should be nil")
	}

	if ec.ThisOrNil() != nil {
		t.Errorf("all errors should be nil")
	}
}
