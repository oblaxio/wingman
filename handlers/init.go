package handlers

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Init handler")
}
