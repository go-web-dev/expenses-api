package controllers

import (
	"net/http"
)

type signupper interface {
	Signup()
}

func signup(service signupper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
