package main

import (
	"dist-grep/utils/rgrep"
	"fmt"
	"os/exec"
	"os/user"
)

//Each test output would be placed in here
var outputFile string = "/go/src/dist-grep/test/test_logs/log"

//Regular Expression Pattern
var regularExpPattern string = "-t -E 1.2{5}\n"
var regularExpPatternFile string = "/go/src/dist-grep/test/golden_logs/regular_pattern_expected_output.log"

// Infrequent Pattern only one match
var infrequentPattern string = "-t -E 1.1{6}\n"
var infrequentPatternFile string = "/go/src/dist-grep/test/golden_logs/infrequent_pattern_expected_output.log"

// Frequent Pattern 71544 matches
var frequentPattern string = "-t -E 2.11\n"
var frequentPatternFile string = "/go/src/dist-grep/test/golden_logs/frequent_pattern_expected_output.log"

// Crash a server by a shell script
var crashShell string = "/go/src/dist-grep/scripts/kill_server_remote.sh"
var crashFile string = "/go/src/dist-grep/test/golden_logs/crashed_pattern_expected_output.log"

func main() {
	// Regular Expression Pattern Test
	r := rgrep.Rgrep{}
	r.Launch(regularExpPattern)

	usr, _ := user.Current()
	usrHome := usr.HomeDir

	// We use diff command in Unix/Linux to compare
	// two logs query output
	cmd := exec.Command("diff", usrHome+outputFile, usrHome+regularExpPatternFile)
	stdOutErr, _ := cmd.CombinedOutput()
	// If there are no output of stdout,
	// the two query results are the same
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Regular Expression Pattern: %s\n", regularExpPattern)
	} else {
		fmt.Printf("Test Failed for Regular Expression Pattern: %s\n", regularExpPattern)
	}

	// Infrequent Pattern Test
	r = rgrep.Rgrep{}
	r.Launch(infrequentPattern)
	cmd = exec.Command("diff", usrHome+outputFile, usrHome+infrequentPatternFile)
	stdOutErr, _ = cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Infrequent Pattern: %s\n", infrequentPattern)
	} else {
		fmt.Printf("Test Failed for Infrequent Pattern: %s\n", infrequentPattern)
	}

	// Frequent Pattern Test
	r = rgrep.Rgrep{}
	r.Launch(frequentPattern)
	cmd = exec.Command("diff", usrHome+outputFile, usrHome+frequentPatternFile)
	stdOutErr, _ = cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Frequent Pattern: %s\n", frequentPattern)
	} else {
		fmt.Printf("Test Failed for Frequent Pattern: %s\n", frequentPattern)
	}

	// We crash a server in VM09, and run frequent pattern test again
	cmd = exec.Command(usrHome+crashShell, "09")
	cmd.Run()
	r = rgrep.Rgrep{}
	r.Launch(frequentPattern)
	cmd = exec.Command("diff", usrHome+outputFile, usrHome+crashFile)
	stdOutErr, _ = cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Crashed One Server Pattern: %s\n", frequentPattern)
	} else {
		fmt.Printf("Test Failed for Crashed One Server Pattern: %s\n", frequentPattern)
	}
}
