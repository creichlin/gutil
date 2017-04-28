package testin

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type process func(source string, operation string, t *testing.T) string

func RunMapTests(t *testing.T, folder string, process process) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		RunMapTest(t, filepath.Join(folder, file.Name()), process)
	}
}

// RunMapTest will read the given file and split it's
// content into a source part and a map of result parts
// It's using delimiters like:
// #########
// # name
// #########
// to parse the result parts.
// For each such part the process func is called
// with the source and the name of the result
// part. The result of the process is the transformation
// from the source part which then is compared with the defined
// result part.
func RunMapTest(t *testing.T, file string, process process) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	content := string(bytes)

	parts, err := splitParts(content)
	if err != nil {
		t.Fatal(err)
	}

	keys := []string{}
	for key := range parts.results {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		t.Run(file+"-"+key, func(t *testing.T) {
			result := process(parts.main, key, t)
			if result != parts.results[key] {
				t.Errorf("Failed to generate expected result for %v:\nexpected:\n%v\nactual:\n%v", key, parts.results[key], result)
			}
		})
	}
}

type parts struct {
	main    string
	results map[string]string
}

func splitParts(content string) (*parts, error) {

	result := &parts{
		results: map[string]string{},
	}
	lastIndex := 0
	nextName := ""

	lines := strings.Split(content, "\n")
	for index := 0; index < len(lines); index++ {
		if name, isLimiter := isLimiter(lines, index); isLimiter {
			if lastIndex == 0 {
				result.main = strings.Join(lines[:index], "\n")
			} else {
				result.results[nextName] = strings.Join(lines[lastIndex:index], "\n")
			}
			lastIndex = index + 3
			nextName = name
		}
	}
	if lastIndex == 0 {
		result.main = content
	} else {
		result.results[nextName] = strings.Join(lines[lastIndex:], "\n")
	}

	return result, nil
}

func isLimiter(lines []string, index int) (string, bool) {
	// we need at least three lines left for an additional part
	if index+2 < len(lines) {
		if strings.HasPrefix(lines[index], "#####") &&
			strings.HasPrefix(lines[index+2], "#####") &&
			strings.HasPrefix(lines[index+1], "#") {
			name := strings.Trim(lines[index+1], " #")
			return name, true
		}
	}
	return "", false
}
