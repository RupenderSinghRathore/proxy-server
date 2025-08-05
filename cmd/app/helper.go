package main

import "net/http"

func makeRequest(method, url string, r *http.Request) (*http.Response, error) {
	req, err := http.NewRequest(method, url, r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res.Header.Add("X-Cache", "MISS")
	return res, err
}
