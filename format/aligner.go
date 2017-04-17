package format

import "strings"

// AlignIndentedBlocks takes a multiline string as input.
// This input is partitioned into blocks where all lines in
// a row having the same indentation belong to the same block.
// Then for each block the Align function is applied.
func AlignIndentedBlocks(input, token string) string {
	lines := strings.Split(input, "\n")
	currentIndentation := countIndentation(lines[0])
	lastIndex := 0
	indentedList := []string{}

	for index, line := range lines {
		indentation := countIndentation(line)
		if currentIndentation != indentation {
			indentedList = append(indentedList, AlignList(lines[lastIndex:index], token)...)
			currentIndentation = indentation
			lastIndex = index
		}
	}

	indentedList = append(indentedList, AlignList(lines[lastIndex:], token)...)

	return strings.Join(indentedList, "\n")
}

func countIndentation(line string) int {
	count := 0
	for _, char := range line {
		if char == ' ' {
			count++
		} else if char == '\t' {
			count += 4
		} else {
			return count
		}
	}
	// this means there is an empty line
	// we would probably better mark that with
	// a second return param or a negative value
	return count
}

// Align will take a multi-line string as input.
// For each line, spaces will be inserted before the first
// such token till all tokens have the same indentation.
func Align(input, token string) string {
	lines := strings.Split(input, "\n")
	return strings.Join(AlignList(lines, token), "\n")
}

// AlignList works like Align, except that input and output
// are lists of strings.
func AlignList(lines []string, token string) []string {
	max := 0

	for _, line := range lines {
		pos := strings.Index(line, token)
		if pos > max {
			max = pos
		}
	}

	modLines := []string{}
	for _, line := range lines {
		pos := strings.Index(line, token)
		if pos != -1 {
			diff := max - pos
			modLines = append(modLines, line[:pos]+strings.Repeat(" ", diff)+line[pos:])
		} else {
			modLines = append(modLines, line)
		}
	}

	return modLines
}
