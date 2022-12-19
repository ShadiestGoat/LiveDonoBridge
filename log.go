package main

import "github.com/fatih/color"

func PrintErr(msg string, items ...any) {
	color.New(color.FgRed).Printf(msg, items...)
}