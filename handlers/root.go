package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func RootHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Root handler")
}
