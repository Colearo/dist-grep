package main

import (
	"dist-grep/utils/rgrep"
	"fmt"
	"os/exec"
	"os/user"
)

var outputFile string = "/go/src/dist-grep/test/test_logs/log"
var regularExpPattern string = "-E 1.2{5}\n"
var regularExpPatternFile string = "/go/src/dist-grep/test/golden_logs/regular_pattern_expected_output.log"
var infrequentPattern string = "-E 1.1{6}\n"
var infrequentPatternFile string = "/go/src/dist-grep/test/golden_logs/infrequent_pattern_expected_output.log"
var frequentPattern string = "-E 2.11\n"
var frequentPatternFile string = "/go/src/dist-grep/test/golden_logs/frequent_pattern_expected_output.log"

func main() {
	r := rgrep.Rgrep{}
	r.Launch(regularPattern)

	usr, _ := user.Current()
	usrHome := usr.HomeDir

	cmd := exec.Command("diff", usrHome+outputFile, usrHome+regularExpPatternFile)
	stdOutErr, _ := cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Regular Expression Pattern: %s\n", regularExpPattern)
	} else {
		fmt.Printf("Test Failed for Regular Expression Pattern: %s\n", regularExpPattern)
		fmt.Println(string(stdOutErr))
	}

	cmd = exec.Command("diff", usrHome+outputFile, usrHome+infrequentPatternFile)
	stdOutErr, _ = cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Infrequent Pattern: %s\n", infrequentPattern)
	} else {
		fmt.Printf("Test Failed for Infrequent Pattern: %s\n", infrequentPattern)
		fmt.Println(string(stdOutErr))
	}

	cmd = exec.Command("diff", usrHome+outputFile, usrHome+frequentPatternFile)
	stdOutErr, _ = cmd.CombinedOutput()
	if len(stdOutErr) == 0 {
		fmt.Printf("Test Passed for Frequent Pattern: %s\n", frequentPattern)
	} else {
		fmt.Printf("Test Failed for Frequent Pattern: %s\n", frequentPattern)
		fmt.Println(string(stdOutErr))
	}

}
