package style

import "fmt"

/*
 *	Couple utility functions that make it easier to add styling
 */

func SetColor(bg string, fg string) string {
	return fmt.Sprintf("%%{B%s}%%{F%s}", bg, fg)
}

func Color(bg string, fg string, text string) string {
	return SetColor(bg, fg) + text + SetColor("-", "-")
}

func Underline(text string) string {
	return fmt.Sprintf("%%{+u}%s%%{-u}", text)
}

func Overline(text string) string {
	return fmt.Sprintf("%%{+o}%s%%{-o}", text)
}

func Button(text string, command string) string {
	return fmt.Sprintf("%%{A:%s:}%s%%{A}", command, text)
}
