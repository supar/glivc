package main

import (
	"gopkg.in/libgit2/git2go.v22"
	"net/url"

	"errors"
	"strconv"
	"strings"
	"time"
)

const TimePattern = "2006-01-02T15:04:05"

type Filter struct {
	Project     *Project
	Branch      *Branch
	Merges      bool
	Authors     []string
	Since, Till time.Time
	Message     string
	Page        *Pages
}

type Pages struct {
	iter        uint
	Page, Limit uint
}

// Create filter struct according to the Query Values
func NewFilter() *Filter {
	return new(Filter)
}

// Set project data
func (this *Filter) SetProject(p string) (err error) {
	var (
		proj int
		resp *Response

		pobj = new(Project)
	)

	proj, err = strconv.Atoi(p)
	if err != nil || proj <= 0 {
		err = errors.New("Project id is not valid")

		return
	}

	resp = get_GitLab("/api/v3/projects/"+strconv.Itoa(proj), pobj)

	if resp.Error != nil {
		log.Error("Unable to open project with id %d: %s", proj, resp.Error.Message)

		err = errors.New(resp.Error.Message)

		return
	}

	this.Project = pobj

	return
}

// Set branch according to the project
func (this *Filter) SetBranch(b string) (err error) {
	var (
		resp *Response

		bobj = new(Branch)
	)

	resp = get_GitLab("/api/v3//projects/"+strconv.Itoa(this.Project.Id)+"/repository/branches/"+b, bobj)

	if resp.Error != nil {
		log.Error("Unable to open branch %s: %s", b, resp.Error.Message)

		err = errors.New(resp.Error.Message)

		return
	}

	this.Branch = bobj

	return
}

// Create branch reference
func (this *Filter) BranchRef() string {
	var (
		branch string
	)

	if this.Branch == nil || this.Branch.Name == "" {
		branch = "master"
	} else {
		branch = this.Branch.Name
	}

	return "refs/heads/" + strings.TrimPrefix(branch, "/")
}

// Parse query arguments
func (this *Filter) ParseQuery(args url.Values) (err error) {
	var (
		systime = time.Now()
	)

	// Check merge flag
	if _, ok := args["merge"]; ok {
		if v, err := strconv.Atoi(args.Get("merge")); err == nil && v > 0 {
			this.Merges = true
		} else {
			this.Merges = false
		}
	} else {
		this.Merges = false
	}

	// Check Author flag
	if _, ok := args["author"]; ok {
		this.Authors = make([]string, 0)

		for _, v := range args["author"] {
			if len(strings.Trim(v, " ")) > 0 {
				this.Authors = append(this.Authors, v)
			}
		}
	}

	// Grep message
	if v := strings.Trim(args.Get("msg"), " "); len(v) > 0 {
		this.Message = v
	}

	// Check period. Since argument
	if v := args.Get("since"); len(v) > 0 {
		this.Since, err = time.ParseInLocation(TimePattern, v, systime.Location())

		if err != nil {
			err = errors.New("Since: " + err.Error())
			return
		}
	}

	// Check period. Till argument
	if v := args.Get("till"); len(v) > 0 {
		this.Till, err = time.ParseInLocation(TimePattern, v, systime.Location())

		if err != nil {
			err = errors.New("Since: " + err.Error())
			return
		}
	}

	// Set page object
	vp := strings.Trim(args.Get("page"), " ")
	vl := strings.Trim(args.Get("limit"), " ")
	if len(vp) > 0 || len(vl) > 0 {
		var p, l int

		p, err = strconv.Atoi(vp)

		if vp != "" && err != nil {
			err = errors.New("Page: " + err.Error())
			return
		}

		if len(vl) > 0 {
			l, err = strconv.Atoi(vl)

			if err != nil {
				err = errors.New("Limit: " + err.Error())
				return
			}
		} else {
			l = 50
		}

		if l > 0 && p == 0 {
			p = 1
		}

		this.Page = &Pages{
			Page:  uint(p),
			Limit: uint(l),
		}
	}

	log.Debug("Filter prepared: %v", this)

	return
}

// Compare signature with filter values
func (this *Filter) IsAuthor(sig *git.Signature) bool {
	if this.Authors == nil || len(this.Authors) == 0 {
		return false
	}

	for _, v := range this.Authors {
		if len(v) > 0 && (strings.Contains(sig.Name, v) || strings.Contains(sig.Email, v)) {
			return true
		}
	}

	return false
}

// Grep message
func (this *Filter) IsMessage(s string) bool {
	str := strings.ToLower(s)
	if len(this.Message) == 0 || strings.Contains(str, strings.ToLower(this.Message)) {
		return true
	}

	return false
}

// Check since argument with sinature
func (this *Filter) IsSince(t time.Time) bool {
	if this.Since.IsZero() || this.Since.Equal(t) || this.Since.Before(t) {
		return true
	}

	return false
}

// Check till argument with signature
func (this *Filter) IsTill(t time.Time) bool {
	if this.Till.IsZero() || this.Till.Equal(t) || this.Till.After(t) {
		return true
	}

	return false
}

// Page counter
func (this *Filter) IsPage() (i int8) {
	i = 0

	if this.Page == nil {
		i = 0
	} else {
		switch {
		case (this.Page.iter / this.Page.Limit) >= this.Page.Page:
			i = 1

		case ((this.Page.iter / this.Page.Limit) + 1) < this.Page.Page:
			i = -1
		}

		this.Page.iter++
	}

	return
}
