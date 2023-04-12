package handlers

import (
	"fmt"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/spf13/cobra"
)

func InitHandler(cmd *cobra.Command, args []string) {
	conf := config.NewConfig()
	if err := conf.Create(); err != nil {
		fmt.Println("Could not initialize wingman config file:", err)
	}
}
