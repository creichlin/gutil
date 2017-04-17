package format

import (
	"github.com/creichlin/gutil/testin"
	"testing"
)

func TestAlign(t *testing.T) {

	testin.RunMapTests(t, "testcases/aligner", func(source, operation string, t *testing.T) string {
		if operation == "aligned" {
			return Align(source, ":")
		}
		if operation == "indented block aligned" {
			return AlignIndentedBlocks(source, ":")
		}
		t.Fatalf("Wrong operation %v", operation)
		return ""
	})
}
