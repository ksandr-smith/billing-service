package apiserver

import (
	"billing-service/internal/app/model"
	"billing-service/internal/app/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/api/v2/wallet/topup", s.handleTopUpBalance()).Methods("POST")
	s.router.HandleFunc("/api/v2/wallet/write_off", s.handleBuy()).Methods("POST")
	s.router.HandleFunc("/api/v2/wallet/{id}", s.handleBalance()).Methods("GET")
	s.router.HandleFunc("/api/v2/wallet/transfer", s.handleTransferBetweenUsers()).Methods("POST")
	s.router.HandleFunc("/api/v2/wallet/{id}/transactions", s.handleGetTransactions()).Methods("GET")
}

func (s *server) handleTopUpBalance() http.HandlerFunc {
	type request struct {
		UserId int `json:"user_id"`
		Amount int `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		wallet, err := s.store.Wallet().FindByUserId(req.UserId)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err = s.store.Wallet().UpdateBalance(wallet, req.Amount); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		tx := &model.Transaction{
			WalletID: wallet.ID,
			Amount:   req.Amount,
			Desc:     "Add User Balance",
		}

		if err = s.store.Transaction().Create(tx); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleBuy() http.HandlerFunc {
	type request struct {
		UserId int `json:"user_id"`
		Amount int `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		wallet, err := s.store.Wallet().FindByUserId(req.UserId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err2 := s.store.Wallet().UpdateBalance(wallet, -req.Amount); err2 != nil {
			s.error(w, r, http.StatusBadRequest, err2)
			return
		}

		tx := &model.Transaction{
			WalletID: wallet.ID,
			Amount:   -req.Amount,
			Desc:     "Purchase",
		}

		if err = s.store.Transaction().Create(tx); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleBalance() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]

		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("user_id is missing in parameters"))
			return
		}
		userId, err := strconv.Atoi(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("user_id must be int"))
			return
		}

		wallet, err := s.store.Wallet().FindByUserId(userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		currency := r.URL.Query().Get("currency")

		if len(currency) == 3 {
			s.logger.Info("currency: " + currency)
			type response struct {
				Success bool               `json:"success"`
				Rates   map[string]float32 `json:"rates"`
			}
			resp := &response{}
			apiUrl := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=0b9e37212bf715e5da9f4e444bb6c4cb")

			respApi, err := http.Get(apiUrl)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			json.NewDecoder(respApi.Body).Decode(resp)
			balance := resp.Rates[currency] * (float32(wallet.Balance) / resp.Rates["RUB"])
			if resp.Success {
				result := fmt.Sprintf("%.2f", balance)
				s.respond(w, r, http.StatusOK, struct {
					Balance  string `json:"balance"`
					Currency string `json:"currency"`
				}{
					Balance:  result,
					Currency: currency,
				})
				return
			}
			s.error(w, r, http.StatusOK, errors.New("api server response is not success"))
			return
		}

		s.respond(w, r, http.StatusOK, wallet)
	}
}

func (s *server) handleTransferBetweenUsers() http.HandlerFunc {
	type request struct {
		SenderId      int `json:"sender_id"`
		DestinationId int `json:"destination_id"`
		Amount        int `json:"amount"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if req.Amount <= 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("amount must be greater 0"))
			return
		}

		sender, err := s.store.Wallet().FindByUserId(req.SenderId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		destination, err2 := s.store.Wallet().FindByUserId(req.SenderId)
		if err2 != nil {
			s.error(w, r, http.StatusInternalServerError, err2)
			return
		}

		if err = s.store.Wallet().UpdateBalance(sender, -req.Amount); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		senderTransaction := &model.Transaction{
			WalletID: req.SenderId,
			Amount:   req.Amount,
			Desc:     "Transfer to user",
		}

		if err = s.store.Transaction().Create(senderTransaction); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		destinationTransaction := &model.Transaction{
			WalletID: req.SenderId,
			Amount:   req.Amount,
			Desc:     "Transfer from user",
		}

		if err = s.store.Transaction().Create(destinationTransaction); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err = s.store.Wallet().UpdateBalance(destination, req.Amount); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetTransactions() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]

		if !ok {
			s.error(w, r, http.StatusBadRequest, errors.New("user_id is missing in parameters"))
			return
		}
		userId, err := strconv.Atoi(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("user_id must be int"))
			return
		}

		wallet, err := s.store.Wallet().FindByUserId(userId)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		var transactions []model.Transaction
		transactions, err = s.store.Transaction().FindByWalletId(wallet.ID, r.URL.Query())
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, transactions)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
