package handler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type walletHandler struct {
	logger *logrus.Logger
}

func NewWalletHandler(logger *logrus.Logger) Handler {
	return &walletHandler{
		logger: logger,
	}
}

func (h *walletHandler) Register(router *httprouter.Router) {
	router.GET("api/v2/wallet/topup", h.TopUpBalance)
}

func (h *walletHandler) TopUpBalance(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	return
}
