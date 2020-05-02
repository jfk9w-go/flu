package http

import "net/http"

type Authorization interface {
	SetAuth(*http.Request)
}

type basicAuth [2]string

func (a basicAuth) SetAuth(req *http.Request) {
	req.SetBasicAuth(a[0], a[1])
}

func Basic(username, password string) Authorization {
	return basicAuth{username, password}
}

type bearerAuth string

func (a bearerAuth) SetAuth(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+string(a))
}

func Bearer(token string) Authorization {
	return bearerAuth(token)
}
