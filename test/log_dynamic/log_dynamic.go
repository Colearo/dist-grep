package main

import (
	"dist-grep/utils/rgrep"
	"fmt"
	"os/exec"
	"os/user"
)

var knownPattern string = "KNOWN PATTERN 1000\n"
var onlyOnePattern string = "[125]\n"
var coverAllPattern string = "Gen logs here\n"

func main() {
	r := rgrep.Rgrep{}
	r.Launch(knownPattern)
	expected_num := 1000 * 10
	test_num := r

	if len(stdOutErr) == expected_num {
		fmt.Printf("Test Passed for Known Pattern: %s\n", knownPattern)
	} else {
		fmt.Printf("Test Failed for Known Pattern: %s\n", knownPattern)
	}

	r.Launch(knownPattern)
	expected_num := 1000 * 10
	test_num := r

	if len(stdOutErr) == expected_num {
		fmt.Printf("Test Passed for Known Pattern: %s\n", knownPattern)
	} else {
		fmt.Printf("Test Failed for Known Pattern: %s\n", knownPattern)
	}

}
