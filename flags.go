package main

import (
	"flag"
	"os"
	"strings"
)

var (
	Host    string
	GitUrl  string
	AToken  string
	GitRoot string
	LogTo   string

	LogFile bool

	DEBUG bool
)

func ParseFlags() {
	flag.StringVar(&Host, "h", "127.0.0.1:8008", "Default host is 127.0.0.1:8008")
	flag.StringVar(&GitUrl, "gu", "https://git.localhost", "Default git URL https://git.localhost")
	flag.StringVar(&AToken, "gt", "", "Private token is used to access application resources without authentication")
	flag.StringVar(&GitRoot, "gr", "/home/git/repositories", "Repository root. Default: /home/git/repositories")
	flag.StringVar(&LogTo, "l", "console", "Write journal to file or console. Default: console")

	flag.BoolVar(&DEBUG, "v", false, "Debug mode. Default false")

	flag.Parse()

	LogFile = false

	if strings.ToLower(LogTo) != "console" {
		if _, err := os.Stat(LogTo); err != nil && os.IsNotExist(err) {
			print("Valid log file required\n")
			os.Exit(-1)
		} else {
			LogFile = true
		}
	}
}
