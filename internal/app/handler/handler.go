package handler

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler interface {
	Register(router *httprouter.Router)
	TopUpBalance(w http.ResponseWriter, r *http.Request, param httprouter.Params)
}
