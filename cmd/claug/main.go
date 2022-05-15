package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/sergei-dyshel/claug/internal/config"

	"github.com/sergei-dyshel/claug/cmd/claug/flags"

	"github.com/spf13/cobra/doc"

	"github.com/invopop/jsonschema"
	"github.com/sergei-dyshel/claug/internal/utils"
	"github.com/sergei-dyshel/claug/pkg/tmux"
	"github.com/spf13/cobra"
)

var configSchemaCmd = &cobra.Command{
	Use: "config-schema",
	Run: func(cmd *cobra.Command, args []string) {
		r := new(jsonschema.Reflector)
		schema := r.Reflect(&config.Config{})
		data, err := json.MarshalIndent(schema, "" /* prefix */, "  " /* indent */)
		utils.AssertErr(err)
		fmt.Println(string(data))
	},
}

// TODO: move where appropriate
func runRoot(cmd *cobra.Command, args []string) error {
	err := tmux.InsertText(tmux.CmdLine, "some text", true)
	if err != nil {
		panic(err)
	}
	return nil
}

func makeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: path.Base(os.Args[0]),
		// RunE:                  runRoot,
		Args: cobra.NoArgs,
		// DisableFlagsInUseLine: true,
	}
	addCommonFlags(cmd)
	prepareCmd(cmd)
	addCommandGroup(cmd, "Server management", startServerCmd(), pingCmd, exitServerCmd)
	addCommandGroup(cmd, "History", pressEnterCmd)
	addCommandGroup(cmd, "Misc", dumpDefaultConfigCmd())
	addCommandGroup(cmd, "Internal", genManCmd(), configSchemaCmd)
	return cmd
}

func genManCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "gen-man",
		Short:  "Generate manpage",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Title:   "MINE",
				Section: "3",
			}
			return doc.GenManTree(cmd.Parent(), header, flags.ManOutDir)
		},
	}
	cmd.Flags().StringVarP(&flags.ManOutDir, "out", "o", ".", "Output manpages to `dir`")
	return cmd
}

func dumpDefaultConfigCmd() *cobra.Command {
	var writeToFile bool

	cmd := &cobra.Command{
		Use:   "default-config",
		Short: "Dump default configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			file := os.Stderr
			if writeToFile {
				configFname, err := getConfigFileName()
				if err != nil {
					return err
				}
				fmt.Printf("Writing to %s\n", configFname)
				file, err = os.Create(configFname)
				if err != nil {
					return utils.Wrapf(err, "could not open %s for writing", configFname)
				}
			}
			_, err := fmt.Fprintln(file, config.Default)
			utils.AssertErr(err)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&writeToFile, "write", "w", false, utils.Doc(`
		Write config to file specified by "--config" option (or its default).
		Othewise write to stdout.
	`))
	return cmd
}

func main() {
	setupCobra()
	if err := makeRootCmd().Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %s", utils.ErrFull(err))
		utils.AssertErr(err)
		os.Exit(1)
	}
}
