package controller

import (
	"net/http"
	"time"
)

func (c *Controller) initHTTPClient() {
	client := &http.Client{
		Timeout: 15 * time.Second, // Global timeout per request
		Transport: &http.Transport{
			MaxIdleConns:          100, // Maximum idle connections across all hosts
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			// You can also tweak DialContext here for faster DNS lookups or timeouts
		},
	}
	c.httpClient = client
}
