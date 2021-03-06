package main

import (
	"github.com/go-martini/martini"
)

var (
	VERSION string = "develop"
	srv     *martini.Martini
	log     *Log
)

func init() {
	// Parse flags
	ParseFlags()

	// Create logger
	log = NewLogger(1000)

	// Start console logger
	if LogFile {
		log.SetLogger("file", `{"filename":"`+LogTo+`"}`)
	} else {
		log.SetLogger("console", "")
	}
	// Set log level
	if DEBUG {
		log.SetLevel(LevelDebug)
	} else {
		log.SetLevel(LevelNotice)
	}

	// Init martini server
	srv = martini.New()

	// Setup middleware
	srv.Use(martini.Recovery())
	srv.Use(LogMiddle())
}

func main() {
	// Close logger
	defer log.Close()
	// Greeting
	log.Info("Dude, i'm starting, be careful, GitLab interface to view commits v-%s", VERSION)

	r := martini.NewRouter()

	r.Get("/", func() string {
		return "Hello world!"
	})

	r.Get("/project", getProjects)
	r.Get("/project/:id", getProjects)
	r.Get("/project/:id/branches", getBranches)
	r.Get("/project/:id/branches/:branch", getBranches)
	r.Get("/project/:id/branch/commits/:branch", getBranchCommits)

	srv.Action(r.Handle)
	srv.RunOnAddr(Host)
}
