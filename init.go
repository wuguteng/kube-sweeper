package main

import (
	"flag"
	"fmt"
	"os"
)

// Variables
var (
	BuildVersion = "v0.0.1"
	BuildName    = "Kubernetes Sweeper"
	BuiltTime    = ""
	CommitID     = ""
)

func init() {
	var showVer bool
	var loggerFile string

	flag.BoolVar(&showVer, "v", false, "Build version")
	flag.StringVar(&loggerFile, "log", "", "Logger file")

	flag.Parse()

	if showVer {
		fmt.Printf("Build name:\t%s\n", BuildName)
		fmt.Printf("Build version:\t%s\n", BuildVersion)
		fmt.Printf("Built time:\t%s\n", BuiltTime)
		fmt.Printf("Commit ID:\t%s\n", CommitID)
		os.Exit(0)
	}
}
