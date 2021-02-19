package main

import (
	"os"

	"github.com/danieloleynyk/dirstatsbeat/cmd"

	_ "github.com/danieloleynyk/dirstatsbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
