package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	client = &http.Client{}
)

// Send http request to the target host
func get_GitLab(url string, data interface{}) (resp *Response) {
	resp = &Response{
		Success: true,
		Results: data,
	}

	req, err := http.NewRequest("GET", GitUrl+url, nil)
	if err != nil {
		resp.SetError(500, err)

		return
	}

	req.Header.Add("PRIVATE-TOKEN", AToken)

	log.Debug("Calling: %s", req.URL.Path)
	res, err := client.Do(req)

	if err != nil {
		resp.SetError(500, err)

		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		resp.SetError(500, err)

		return
	} else {
		resp.SetData(body)
	}

	if res.StatusCode != 200 {
		resp.SetError(res.StatusCode, res.Status)

		return
	}

	// Try to decode recevied data to the result section
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	err = decoder.Decode(&resp.Results)
	if err != nil {
		log.Warn(err.Error())
	}

	return
}
