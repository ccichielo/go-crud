package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountByID))
	router.HandleFunc("/account/{from}/{to}/{amount}", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API Server running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccount(w)
	case "POST":
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	switch method := r.Method; method {
	case "GET":
		return s.handleGetAccountByID(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Invalid account ID format")
	}

	account, err := s.store.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("Account not found")
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	request := &CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return err
	}

	account := NewAccount(request.FirstName, request.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("Invalid account ID format")
	}

	err = s.store.DeleteAccount(accountID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accountID)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	from := mux.Vars(r)["from"]
	to := mux.Vars(r)["to"]
	amount := mux.Vars(r)["amount"]

	fromID, err := strconv.Atoi(from)
	if err != nil {
		return fmt.Errorf("Invalid account ID format for 'from'")
	}

	toID, err := strconv.Atoi(to)
	if err != nil {
		return fmt.Errorf("Invalid account ID format for 'to'")
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return fmt.Errorf("Invalid format for 'amount'")
	}

	err = s.store.Transfer(fromID, toID, amountInt)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Successfully transferred %d from Account %d to Account %d", amountInt, fromID, toID)
	return WriteJSON(w, http.StatusOK, msg)
}

type APIError struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
