package main

import (
	"fmt"
	"log"
	"os"

	ct "github.com/daviddengcn/go-colortext"
)

type Logger struct {
	ErrorIcon   string
	SuccessIcon string
	WarningIcon string
	LoadingIcon string
	LogIcon     string
}

func (logger Logger) Error(message string, params ...any) {
	ct.Foreground(ct.Red, true)
	log.Printf(logger.ErrorIcon+"  "+message, params...)
	ct.ResetColor()
	fmt.Println("\npress any key to continue...")
	Scanner.Scan()
	logger.Loading("closing application")
	os.Exit(1)
}
func (logger Logger) Loading(message string, params ...any) {
	ct.Foreground(ct.Magenta, true)
	log.Printf(logger.LoadingIcon+"  "+message+"...", params...)
	ct.ResetColor()
}

func (logger Logger) Success(message string, params ...any) {
	ct.Foreground(ct.Green, true)
	log.Printf(logger.SuccessIcon+"  "+message, params...)
	ct.ResetColor()
}

func (logger Logger) Warning(message string, params ...any) {
	ct.Foreground(ct.Yellow, true)
	log.Printf(logger.WarningIcon+"  "+message, params...)
	ct.ResetColor()
}

func (logger Logger) Log(message string, params ...any) {
	ct.ResetColor()
	log.Printf(logger.LogIcon+"  "+message, params...)
}

var Console = Logger{
	ErrorIcon:   "💥",
	SuccessIcon: "✅",
	WarningIcon: "⚠️",
	LoadingIcon: "⌛",
	LogIcon:     "ℹ️",
}
