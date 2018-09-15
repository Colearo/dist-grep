package main

import "dist-grep/utils/rgrep"

func main() {
	r := rgrep.Rgrep {}
	r.Launch("-e 1.22222\n")
}
