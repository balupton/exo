package main

import (
	"os"

	"github.com/deref/exo/config"
	"github.com/deref/exo/import/procfile"
	"github.com/deref/exo/util/cmdutil"
)

func main() {
	cfg, err := procfile.Import(os.Stdin)
	if err != nil {
		cmdutil.Fatal(err)
	}
	if err := config.Generate(os.Stdout, cfg); err != nil {
		cmdutil.Fatal(err)
	}
}
