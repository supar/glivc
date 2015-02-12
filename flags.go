package main

import "flag"

var (
	Host		string
	GitUrl		string
	AToken		string
	GitRoot		string
	
	DEBUG		bool
)

func ParseFlags() {
	flag.StringVar(&Host, "h", "127.0.0.1:8008", "Default host is 127.0.0.1:8008")
	flag.StringVar(&GitUrl, "gu", "https://git.localhost", "Default git URL https://git.localhost")
	flag.StringVar(&AToken, "gt", "", "Private token is used to access application resources without authentication")
	flag.StringVar(&GitRoot, "gr", "/home/git/repositories", "Repository root. Default: /home/git/repositories")
	
	flag.BoolVar(&DEBUG, "v", false, "Debug mode. Default false")
	
	flag.Parse()
}