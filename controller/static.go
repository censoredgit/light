package controller

import "net/http"

type static struct {
	uri     string
	handler http.Handler
}
