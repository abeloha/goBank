package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/transfer", withJwtToken(s.handleTransfer))

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

	token, err := generateToken(account)

	if err != nil {
		log.Printf("Error generating token: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]any{"token": token, "account": account})
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
	if r.Method != "POST" {
		return fmt.Errorf("Unsupported method %s", r.Method)
	}

	transferReq := TransferRequestModel{}
	err := json.NewDecoder(r.Body).Decode(&transferReq)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return fmt.Errorf("Error parsing request")
	}

	if transferReq.Amount < 100 {
		return fmt.Errorf("Amount cannot be less than 100")
	}

	return WriteJSON(w, http.StatusOK, transferReq)

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

func validateJwt(tokenString string) (*jwt.Token, error) {
	mySigningKey := []byte(os.Getenv(os.Getenv("JWT_SECRET")))
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKey, nil
	})
}
func withJwtToken(f apiFuc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Authorization token missing"})
			return
		}

		token, err := validateJwt(tokenString)
		if err != nil {
			log.Println("Error validating JWT token: ", err)
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Error validating request"})
			return
		}

		log.Printf("%+v", token.Claims)

		f(w, r)
	}
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

func generateToken(account *Account) (string, error) {

	mySigningKey := []byte(os.Getenv(os.Getenv("JWT_SECRET")))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":             account.ID,
		"account_number": account.Number,
		"nbf":            time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(mySigningKey)
}
