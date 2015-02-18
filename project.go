package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"fmt"
	"strconv"
)

type Project struct {
	Id		int		`json:"id"`
	Name	string	`json:"name"`
	Path	string	`json:"path_with_namespace"`
}

type Branch struct {
	Name	string	`json:"name"`
}

func getProjects(params martini.Params) (int, string) {
	var (
		resp		*Response
		r_pr		interface{}
		p_id		int
		url = "/api/v3/projects"
	)
	
	if _, ok := params["id"]; ok {
		if id, err := strconv.Atoi(params["id"]); err == nil && id > 0 {
			p_id = id
		}
	}
	
	if p_id > 0 {
		url = fmt.Sprintf("%s/%d", url, p_id)
		r_pr = &Project{}
	} else {
		r_pr = &[]Project{}
	}
	
	resp = get_GitLab(url, r_pr)
	
	if resp.Error != nil {
		log.Error(resp.Error.Message)
	}
	
	return resp.Compile()
}


// Get branches from GitLab API
func getBranches(params martini.Params) (int, string) {
	var (
		resp		*Response
		r_br		interface{}
		p_id		int
		url = "/api/v3/projects/%d/repository/branches"
	)
	
	if _, ok := params["id"]; ok {
		if id, err := strconv.Atoi(params["id"]); err == nil && id > 0 {
			p_id = id
		}
	}
	
	if p_id <= 0 {
		resp.SetError(http.StatusBadRequest,  "Project id is required")
		return resp.Compile()
	}
	
	url = fmt.Sprintf(url, p_id)
	
	if v, ok := params["branch"]; ok && v != "" {
		url = fmt.Sprintf("%s/%s", url, v)
		
		r_br = &Branch{}
	} else {
		r_br = &[]Branch{}
	}
	
	resp = get_GitLab(url, r_br)
	
	if resp.Error != nil {
		log.Error(resp.Error.Message)
	}
	
	return resp.Compile()
}