package main

import "github.com/gszhangwei/open-spdd/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
