//go:build ignore
// +build ignore

package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	prefixLine = "```golang"
	suffixLine = "```"
	codeStart  = regexp.MustCompile(`^[\s]*<!--[\s]*include:.*-->[\s]*$`)
	codeEnd    = regexp.MustCompile(`^[\s]*<!--[\s]*/[\s]*-->[\s]*$`)
	nodocEx    = regexp.MustCompile(`^.*//[\s]*nodoc[+]?[0-9]*[\s]*$`)
)

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseInclude(line string) string {
	rhs := strings.TrimPrefix(strings.TrimPrefix(line, "<!--"), "include:")
	return strings.TrimSpace(strings.TrimSuffix(rhs, "-->"))
}

func nodoc(lines []string) []string {
	var out []string
	for i := 0; i < len(lines); i++ {
		if match := nodocEx.FindString(lines[i]); len(match) > 0 {
			parts := strings.Split(match, "nodoc+")
			if len(parts) == 2 {
				skip, err := strconv.Atoi(parts[1])
				exitOnError(err)
				i += skip
			}
			continue
		}
		out = append(out, lines[i])
	}
	return out
}

func main() {
	file := flag.String("file", "README.md", "source and target file name")
	flag.Parse()

	text, err := os.ReadFile(*file)
	exitOnError(err)

	input := strings.Split(string(text), "\n")

	start := -1
	end := -1
	for i, line := range input {
		switch {
		case codeStart.MatchString(line):
			start = i
		case codeEnd.MatchString(line):
			end = i
		}
	}

	if start < 0 || end < 0 || end < start {
		exitOnError(errors.New("failed to find code include"))
	}

	inc, err := os.ReadFile(parseInclude(input[start]))
	exitOnError(err)
	code := strings.Split(string(inc), "\n")

	out := make([]string, 0, len(input))
	out = append(out, input[:start+1]...)
	out = append(out, prefixLine)
	out = append(out, nodoc(code)...)
	out = append(out, suffixLine)
	out = append(out, input[end:]...)

	output := strings.Join(out, "\n")
	// println(output)
	err = os.WriteFile(*file, []byte(output), 0644)
	exitOnError(err)
}
