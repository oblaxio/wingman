package print

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func Info(msg string) {
	fmt.Printf("%s %s\n", cyan("[INFO]"), bold(msg))
}

func Rebuild(msg string) {
	fmt.Printf("%s %s\n", yellow("[REBUILD]"), bold(msg))
}

func SvcOut(service string, msg string) {
	fmt.Printf("%s\t %s\n", green("["+service+"]"), msg)
}

func SvcErr(service string, msg string) {
	fmt.Printf("%s\t %s\n", red("["+service+"]"), msg)
}
