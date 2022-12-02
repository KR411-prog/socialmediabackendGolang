package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		UserEmail string `json:"userEmail"`
		Text      string `json:"text"`
	}
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	post, err := apiCfg.dbClient.CreatePost(params.UserEmail, params.Text)
	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, post)
	log.Println("post created", post)
	return

}

func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/posts/")

	if id == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("no uuid provided to handlerDeletePost"))
		return
	}

	err := apiCfg.dbClient.DeletePost(id)

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusOK, struct{}{})
	return
}


func (apiCfg apiConfig) handlerRetrievePosts(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("userEmail")

	if userEmail == "" {
	  respondWithError(w, http.StatusBadRequest, errors.New("no userEmail provided to handlerRetrievePosts"))
	  return
  }

	posts,err := apiCfg.dbClient.GetPosts(userEmail)

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusBadRequest, err)
		return
	  }
	  respondWithJSON(w, http.StatusOK,posts)
	  return

}
