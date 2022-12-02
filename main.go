package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KR411-prog/socialmedia/internal/database"
)

const addr = "localhost:8082"

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	if payload != nil {
		response, err := json.Marshal(payload)
		if err != nil {
			log.Println("error marshalling", err)
			w.WriteHeader(500)
			response, _ := json.Marshal(errorBody{
				Error: "error marshalling",
			})
			w.Write(response)
			return
		}
		w.WriteHeader(code)
		w.Write(response)
		return
	}
}

type errorBody struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	if err == nil {
		log.Println("dont call respondwithError with a nil error!")
		return
	}
	log.Println(err)
	respondWithJSON(w, 500, errorBody{
		Error: err.Error(),
	})
	return
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 500, errors.New("server error"))
}

type apiConfig struct {
	dbClient database.Client
}

func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerGetUser(w, r)
	case http.MethodPost:
		apiCfg.handlerCreateUser(w, r)
	case http.MethodPut:
		// call PUT handler
	case http.MethodDelete:
		apiCfg.handlerDeleteUser(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}

}

func (apiCfg apiConfig) endpointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerRetrievePosts(w, r)
	case http.MethodPost:
		apiCfg.handlerCreatePost(w, r)
	case http.MethodPut:
		// call PUT handler
	case http.MethodDelete:
		apiCfg.handlerDeletePost(w, r)
	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}

}

func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	user, err := apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
	log.Println("user created", user)
	return
}

func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")

	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to handlerDeleteUser"))
		return
	}

	err := apiCfg.dbClient.DeleteUser(email)

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
	return
}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")

	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to handlerDeleteUser"))
		return
	}

	user, err := apiCfg.dbClient.GetUser(email)

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
	return

}

func userIsEligible(email, password string, age int) error {
	if email == "" {
		err := errors.New("email can't be empty")
		return err
	}
	if password == "" {
		err := errors.New("email can't be empty")
		return err
	}
	if age < 18 {
		err := fmt.Errorf("age must be at least %d years old", 18)
		return err
	}
	return nil
}

func main() {
	fmt.Println("Welcome to Social Media Backend in Golang")
	const dbPath = "./internal/database/db.json"
	dbClient := database.NewClient(dbPath)

	apiCfg := apiConfig{
		dbClient: *dbClient,
	}

	serveMux := http.NewServeMux()

	srv := http.Server{
		Handler:      serveMux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	serveMux.HandleFunc("/", testHandler)
	serveMux.HandleFunc("/err", testErrHandler)

	serveMux.HandleFunc("/users", apiCfg.endpointUsersHandler)
	serveMux.HandleFunc("/users/", apiCfg.endpointUsersHandler)
	serveMux.HandleFunc("/posts", apiCfg.endpointPostsHandler)
	serveMux.HandleFunc("/posts/", apiCfg.endpointPostsHandler)
	srv.ListenAndServe()
}
