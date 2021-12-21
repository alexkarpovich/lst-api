package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type BaseHanlder struct {
	router *mux.Router
}

func (i *BaseHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.router.ServeHTTP(w, r)
}
