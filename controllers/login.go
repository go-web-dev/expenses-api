package controllers

import (
	"net/http"
)

type loginner interface {
	Login()
}

func login(service loginner) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
