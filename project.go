package main

import (
	"github.com/go-martini/martini"
	
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
		resp *Response
		url = "/api/v3/projects"
	)
	
	if _, ok := params["id"]; ok {
		if id, err := strconv.Atoi(params["id"]); err == nil && id > 0 {
			url = url + "/" + strconv.Itoa(id)
		}
	}
	
	resp = get_GitLab(url, nil)
	
	if resp.Error != nil {
		log.Error(resp.Error.Message)
	}
	
	return resp.Compile()
}