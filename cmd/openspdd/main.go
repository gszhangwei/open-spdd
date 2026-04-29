package main

import "github.com/gszhangwei/open-spdd/cmd"

var version = "dev"

func main() {
	cmd.SetVersion(cmd.ResolveVersion(version))
	cmd.Execute()
}
