package http

import (
	"net/http"
	"net/http/cookiejar"
)

type jwClient struct {
	Username   string
	Password   string
	httpClient *http.Client
}

func NewSpider(username, password string) *jwClient {
	// Allocate a new cookie jar to mimic the browser behavior:
	cookieJar, _ := cookiejar.New(nil)

	// Fill up basic data:
	c := &jwClient{
		Username: username,
		Password: password,
	}

	// When initializing the http.Client, copy default values from http.DefaultClient
	// Pass a pointer to the cookie jar that was created earlier:
	c.httpClient = &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           cookieJar,
		Timeout:       http.DefaultClient.Timeout,
	}
	return c

}
