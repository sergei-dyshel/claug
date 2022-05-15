package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/pflag"

	"github.com/fatih/color"
	"github.com/sergei-dyshel/claug/internal/utils"
	"github.com/spf13/cobra"
)

const (
	miscGroupName   = "Misc commands"
	groupAnnotation = "group"
)

func setupCobra() {
	cobra.EnableCommandSorting = false
	cobra.AddTemplateFunc("wrappedLocalFlagUsages", wrappedLocalFlagUsages)
	cobra.AddTemplateFunc("wrappedInheritedFlagUsages", wrappedInheritedFlagUsages)
	cobra.AddTemplateFunc("commandGroupsHelp", commandGroupsHelp)
}

func wrappedFlagUsages(flags *pflag.FlagSet) string {
	width := 80
	if utils.IsTerminal(os.Stderr) {
		width = utils.TerminalWidth(os.Stderr)
	}
	return flags.FlagUsagesWrapped(width - 1)
}

func wrappedLocalFlagUsages(cmd *cobra.Command) string {
	return wrappedFlagUsages(cmd.LocalFlags())
}

func wrappedInheritedFlagUsages(cmd *cobra.Command) string {
	return wrappedFlagUsages(cmd.InheritedFlags())
}

// rpad adds padding to the right of a string.
// (copied from cobra code)
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

func groupHelp(name string, commands []*cobra.Command) string {
	return "  " + name + "\n" + strings.Join(
		utils.Map(commands, func(cmd *cobra.Command) string {
			return fmt.Sprintf("     %s %s",
				rpad(cmd.Name(), cmd.NamePadding()),
				cmd.Short,
			)
		}), "\n",
	)
}

func commandGroupsHelp(cmd *cobra.Command) string {
	commands := utils.Filter(cmd.Commands(), func(c *cobra.Command) bool {
		return c.IsAvailableCommand() || c.Name() == "help"
	})
	groupName := func(cmd *cobra.Command) string {
		return utils.FirstNonZero(cmd.Annotations[groupAnnotation], miscGroupName)
	}
	groups := utils.Group(commands, groupName)
	return strings.Join(utils.Map(groups, func(group []*cobra.Command) string {
		name := groupName(group[0])
		return groupHelp(name, group)
	}), "\n\n")
}

func prepareCmd(cmd *cobra.Command) {
	tmpl := cmd.UsageTemplate()
	tmpl = strings.ReplaceAll(tmpl, ".LocalFlags.FlagUsages", "wrappedLocalFlagUsages .")
	tmpl = strings.ReplaceAll(tmpl, ".InheritedFlags.FlagUsages", "wrappedInheritedFlagUsages .")
	tmpl = regexp.MustCompile(`(?s)(Available Commands:)(.*?)({{end}}{{if .HasAvailableLocalFlags}})`).
		ReplaceAllString(tmpl, "${1}\n\n{{commandGroupsHelp .}}${3}")
	if utils.IsTerminal(os.Stdout) {
		pat := regexp.MustCompile(`(?m)^[\w ]+:`)
		bold := color.New(color.FgWhite, color.Bold)
		tmpl = pat.ReplaceAllStringFunc(tmpl, func(match string) string {
			return bold.Sprint(match)
		})
	}
	// TEMP: fmt.Println(tmpl)
	cmd.SetUsageTemplate(tmpl)
}

func addCommandGroup(cmd *cobra.Command, name string, subcmds ...*cobra.Command) {
	for _, subcmd := range subcmds {
		if subcmd.Annotations == nil {
			subcmd.Annotations = make(map[string]string)
		}
		subcmd.Annotations[groupAnnotation] = name
		cmd.AddCommand(subcmd)
	}
}
