package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Command groups for organized help display
const (
	BasicCommands   = "Basic Commands"
	VersionCommands = "Version Commands"
)

func NewCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gvm",
		Short: "A tool to manage Go versions",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	InitFlags(rootCmd)

	// Set custom help template to display commands in groups
	rootCmd.SetHelpTemplate(customHelpTemplate)
	rootCmd.SetUsageTemplate(customUsageTemplate)

	// Register custom template function
	cobra.AddTemplateFunc("customCommands", func(cmd *cobra.Command) string {
		return customCommands(cmd)
	})
	return rootCmd
}

// customHelpTemplate is the custom help template that groups commands
var customHelpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// customUsageTemplate provides custom usage formatting with command groups
var customUsageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{. | customCommands}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// Commands returns the grouped commands for help display
func customCommands(cmd *cobra.Command) string {
	groups := make(map[string][]*cobra.Command)

	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		group := getCommandGroup(c)
		groups[group] = append(groups[group], c)
	}

	// Define group order
	groupOrder := []string{BasicCommands, VersionCommands}

	var result strings.Builder
	for _, groupName := range groupOrder {
		if commands, exists := groups[groupName]; exists {
			sort.Slice(commands, func(i, j int) bool {
				return commands[i].Name() < commands[j].Name()
			})

			result.WriteString(fmt.Sprintf("\n%s:\n", groupName))
			for _, c := range commands {
				result.WriteString(fmt.Sprintf("  %-15s %s\n", c.Name(), c.Short))
			}
		}
	}

	// Add any remaining commands that don't have a group
	for groupName, commands := range groups {
		found := false
		for _, orderedGroup := range groupOrder {
			if groupName == orderedGroup {
				found = true
				break
			}
		}
		if !found {
			sort.Slice(commands, func(i, j int) bool {
				return commands[i].Name() < commands[j].Name()
			})
			result.WriteString(fmt.Sprintf("\n%s:\n", groupName))
			for _, c := range commands {
				result.WriteString(fmt.Sprintf("  %-15s %s\n", c.Name(), c.Short))
			}
		}
	}

	return result.String()
}

// getCommandGroup determines which group a command belongs to
func getCommandGroup(cmd *cobra.Command) string {
	// Check if command has an annotation for its group
	if group, exists := cmd.Annotations["group"]; exists {
		return group
	}

	return BasicCommands
}
