package main

import (
	"fmt"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Uhh.. Oh... somethign went wront while executing wingman '%s'", err)
		os.Exit(1)
	}
}
