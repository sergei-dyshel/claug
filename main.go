package main

import (
	"fmt"
	"os"

	"github.com/sergei-dyshel/claug/tmux"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	executor := tmux.CmdLineExecutor{}
	_, err := executor.Run(&tmux.SendKeys{Keys: []string{"Enter"}})
	if err != nil {
		panic(err)
	}
	err = tmux.InsertText(&executor, "some text", true)
	if err != nil {
		panic(err)
	}
	return nil
}

var rootCmd = &cobra.Command{
	RunE:                  run,
	DisableFlagsInUseLine: true,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
