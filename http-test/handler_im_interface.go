package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"net/http"
)

const im_service_url = "http://imnew.appcloud.ztecs/zteim4ios"

func imInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();

	// check parameters
	if (r.Form["method"] == nil) {
		writeResponseMessage(w, 400, "Parameter [method] is required.")
		return
	}
	if (r.Form["jsonRequest"] == nil) {
		writeResponseMessage(w, 400, "Parameter [jsonRequest] is required.")
		return	
	}

	method := r.Form["method"][0]
	jsonRequest := r.Form["jsonRequest"][0]
	fmt.Printf("Requesting im server: [%s] %s \n", method, jsonRequest)

	// request im server
	resp, err := http.PostForm(im_service_url,
			url.Values{"method": { method }, "jsonRequest": { jsonRequest }});

	if (err != nil) {
		writeResponseMessage(w, 500, "Failed to send request to im server.")
		return
	}

	if (resp.StatusCode != 200) {
		writeResponseMessage(w, 500, "Failed to process request: " + resp.Status)
		return
	}

	// read response from im server
	body, err := ioutil.ReadAll(resp.Body)
	if (err != nil) {
		writeResponseMessage(w, 500, "Failed to read response from im server.")
		return
	}

	respStr := string(body[:])
	fmt.Printf("Response from im server: [%s] %s\n", method, respStr)

	// send response
	fmt.Fprint(w, respStr)
}