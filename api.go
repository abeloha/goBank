package main

import (
	"encoding/json"
	"fmt"
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

	router.HandleFunc("/account", makeHttpHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHttpHandlerFunc(s.handleAccountById))

	fmt.Println("Searving API server on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET", "":
		return s.handleGetAccounts(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("Unsupported method %s", r.Method)
	}
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET", "":
		return s.handleGetAccountById(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("Unsupported method %s", r.Method)
	}
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts(200)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)

	if err != nil {
		return fmt.Errorf("Record not found for the given id: %d", id)
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	accountReq := CreateAccountRequestModel{}

	err := json.NewDecoder(r.Body).Decode(&accountReq)

	if err != nil {
		fmt.Println("Error personing json:", err)
		return fmt.Errorf("Error parsing json: %v", err)
	}

	account, err := s.store.CreateAccount(accountReq)
	if err != nil {
		fmt.Println("Error creating account:", err)
		return fmt.Errorf("Error creating account: %v", err)
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		fmt.Println("Error deleting account:", err)
		return fmt.Errorf("Error deleting account. Try again")
	}

	return WriteJSON(w, http.StatusOK, "Record Deleted")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type apiFuc func(http.ResponseWriter, *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
	return nil
}

type ApiError struct {
	Error string `json:"error"`
}

func makeHttpHandlerFunc(f apiFuc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	isStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(isStr)
	if err != nil {
		return 0, fmt.Errorf("Id is not a valid number")
	}

	return id, nil
}
