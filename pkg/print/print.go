package print

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	cyan  = color.New(color.FgCyan).SprintFunc()
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
	bold  = color.New(color.Bold).SprintFunc()
)

func PrintInfo(msg string) {
	fmt.Printf("%s %s\n", cyan("[INFO]"), bold(msg))
}

func PrintSvcOut(service string, msg string) {
	fmt.Printf("%s\t %s\n", green("["+service+"]"), msg)
}

func PrintSvcErr(service string, msg string) {
	fmt.Printf("%s\t %s\n", red("["+service+"]"), msg)
}
