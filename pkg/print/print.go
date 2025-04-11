package print

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const labelLength = 20

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

func Debug(msg string) {
	now := time.Now() // current local time
	sec := now.Unix()
	printer("🙃", "DEBUG", yellow, "["+strconv.Itoa(int(sec))+"] "+msg)
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

func printer(icon string, service string, colorfn func(a ...any) string, msg string) {
	fmt.Printf(
		"%s %s %-"+strconv.Itoa(labelLength)+"s %s %s\n",
		icon,
		colorfn(time.Now().Format(time.TimeOnly)),
		colorfn(strings.ToLower(adjustLabel(service, labelLength))),
		colorfn("│"),
		msg,
	)
}

func adjustLabel(label string, max int) string {
	if len(label) > max {
		return label[:max]
	} else if len(label) < max {
		return label + strings.Repeat(" ", max-len(label))
	}
	return label
}
