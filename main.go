package main

import (
	"os"

	"github.com/Qolzam/ofc-bootstrap/cmd"

	"github.com/Qolzam/ofc-bootstrap/version"
)

func main() {

	if err := cmd.Execute(version.Version, version.GitCommit); err != nil {
		// fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	return
}
