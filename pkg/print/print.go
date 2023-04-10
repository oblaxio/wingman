package print

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	cyan    = color.New(color.FgHiCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	hiblue  = color.New(color.FgHiBlue).SprintFunc()
	white   = color.New(color.FgHiWhite).SprintFunc()
	gray    = color.New(color.FgWhite).SprintFunc()
)

func Info(msg string) {
	printer("👻", "info", magenta, msg)
}

func Rebuild(msg string) {
	printer("🔄", "rebuild", cyan, msg)
}

func SvcOut(service string, msg string) {
	printer("🔹", service, blue, msg)
}

func SvcErr(service string, msg string) {
	printer("🧨", service, red, msg)
}

func SvcWarn(service string, msg string) {
	printer("🟡", service, yellow, msg)
}

func SvcProxy(msg string) {
	printer("🚀", "proxy", white, msg)
}

func printer(icon string, service string, colorfn func(a ...interface{}) string, msg string) {
	fmt.Printf(
		"%s %-18s %s   %s\n",
		icon,
		colorfn(strings.ToUpper(service)),
		gray(time.Now().Format(time.DateTime)),
		msg,
	)
}
