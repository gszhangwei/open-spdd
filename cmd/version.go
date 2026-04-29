package cmd

import "runtime/debug"

var readBuildInfo = debug.ReadBuildInfo

func ResolveVersion(injected string) string {
	if injected != "" && injected != "dev" {
		return injected
	}
	if bi, ok := readBuildInfo(); ok && bi != nil {
		if v := bi.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return injected
}

func SetReadBuildInfoForTest(f func() (*debug.BuildInfo, bool)) (restore func()) {
	prev := readBuildInfo
	readBuildInfo = f
	return func() { readBuildInfo = prev }
}
