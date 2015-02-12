package main

import (
	"gopkg.in/libgit2/git2go.v22"
	"github.com/go-martini/martini"
	"net/http"
	
	"errors"
	"strings"
	"time"
)

type Commit struct {
	Id				string		`json:"id"`
	AuthorName		string		`json:"authorname"`
	AuthorEmail		string		`json:"authoremail"`
	Date			time.Time	`json:"date"`
	Message			string		`json:"message"`
}


// Common function to open repository and create walker
// Need repository name
func getWalker(repName string) (walker *git.RevWalk, err error) {
	var (
		repPath string
	)
	
	if repName == "" {
		return nil, errors.New("Empty repository")
	}
	
	// Remove leading slash
	repPath = GitRoot + "/" + strings.TrimPrefix(repName, "/") + ".git"
	
	// Create repository object
	repo, err := git.OpenRepository(repPath)
	if err != nil {
		return nil, err
	}
	
	// Get repository walker
	return repo.Walk()
}


// Get branch commits
func getBranchCommits(params martini.Params, req *http.Request) (int, string) {
	var (
		err		error
		commits	[]Commit
		resp = &Response{
			Success: true,
			Results: &commits,
		}
		
		filter = NewFilter()
	)
	
	err = filter.SetProject(params["id"])
	if err != nil {
		resp.SetError(http.StatusBadRequest,  "Project not available")
		return resp.Compile()
	}
	
	err = filter.SetBranch(params["branch"])
	if err != nil {
		resp.SetError(http.StatusBadRequest,  "Branch " + params["branch"] + " is not available")
		return resp.Compile()
	}
	
	err = filter.ParseQuery(req.URL.Query())
	
	// Check filter strict params
	if err != nil {
		resp.SetError(http.StatusBadRequest,  err.Error())
		return resp.Compile()
	}
	
	var (
		commitCallback = makeRevCallback(&commits, filter)
	)
	
	// Get repository walker
	walker, err := getWalker(filter.Project.Path)
	
	if err != nil {
		log.Error(err.Error())
		
		resp.SetError(500, "Unable to open repository")
		return resp.Compile()
	}
	
	// Set sorting
	walker.Sorting(git.SortTime ^ git.SortReverse)
	
	// set ref
	err = walker.PushRef(filter.BranchRef())
	
	if err != nil {
		log.Error(err.Error())
		
		resp.SetError(500, "Unable to set refs/heads/master")
		return resp.Compile()
	}
	
	
	err = walker.Iterate(commitCallback)
	if err != nil {
		log.Error(err.Error())
		
		resp.SetError(500, "Unable to walk through commits")
		return resp.Compile()
	}
	
	return resp.Compile()
}


// Make revision iterator callback
func makeRevCallback(ret *[]Commit, filter *Filter) func(commit *git.Commit) bool {
	
	return func(commit *git.Commit) bool {
		var (
			c = Commit{
				Id: commit.Id().String(),
				Message: commit.Message(),
				AuthorName: commit.Author().Name,
				AuthorEmail: commit.Author().Email,
				Date: commit.Author().When,
			}
		)
		
		// Count parent
		if !filter.Merges && commit.ParentCount() > 1 {
			return true
		}
		
		// Check Author
		if filter.Authors != nil && filter.IsAuthor(commit.Author()) == false {
			return true
		}
		
		// Grep message
		if filter.IsMessage(commit.Message()) == false {
			return true
		}
		
		// Check since date
		if filter.IsSince(c.Date) == false || filter.IsTill(c.Date) == false {
			return true
		}
		
		switch filter.IsPage() {
			case -1:
				return true
			
			case 1:
				return false
			
			default:
				*ret = append(*ret, c)
		}
		
		return true
	}
}
