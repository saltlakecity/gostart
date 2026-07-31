package styles

import "os"

const (
	reset = "\033[0m"
	bold  = "\033[1m"

	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

const (
	PromptIcon  = "\u25c6"
	SuccessIcon = "\u2713"
	ErrorIcon   = "\u2717"
)

func enabled() bool {
	_, noColor := os.LookupEnv("NO_COLOR")

	return !noColor && os.Getenv("TERM") != "dumb"
}

func paint(code, text string) string {
	if !enabled() {
		return text
	}

	return code + text + reset
}

func Bold(text string) string {
	return paint(bold, text)
}

func Accent(text string) string {
	return paint(cyan, text)
}

func Success(text string) string {
	return paint(green, text)
}

func Warning(text string) string {
	return paint(yellow, text)
}

func Error(text string) string {
	return paint(red, text)
}

func Muted(text string) string {
	return paint(gray, text)
}
