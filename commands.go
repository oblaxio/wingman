package main

import (
	"github.com/oblaxio/wingman/handlers"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wingman",
	Short: "wingman - a simple CLI to run and restart golang services",
	Long:  "",
	Run:   handlers.RootHandler,
}

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initializes a new wingman project",
	Run:     handlers.InitHandler,
}

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "Runs the wingman project from the config file",
	Args:    cobra.ExactArgs(1),
	Run:     handlers.StartHandler,
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
}
