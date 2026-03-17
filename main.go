/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package main

import "github.com/aae42/propel/cmd"

// Version information set by GoReleaser via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
