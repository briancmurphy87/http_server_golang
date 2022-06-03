package main

import (
	// "fmt"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/briancmurphy87/http_server_golang/internal/database"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(200)
	// w.Write([]byte("{}"))
	// you can use any compatible type, but let's use our database package's User type for practice
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 200, errors.New("test error"))
}

/*
An endpoint would call this function once per request, 
	right before returning. 
It will 
	write some standard HTTP headers, 
	marshal an interface into JSON bytes, 
	and write all that to the response along with a status code.
*/
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	log.Println("IN[respondWithJSON]")
	// TODO: 
	// add headers first 
	// w.Header().Set(key, value)
	w.Header().Set("Content-Type", "application/json")

	// write json body 
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(dat)

	// write status code
	w.WriteHeader(code)

}


type errorBody struct {
	Error string `json:"error"`
}

// This function should take the err, create a new errorBody, then call respondWithJSON.
func respondWithError(w http.ResponseWriter, code int, err error) {
	log.Println("IN[respondWithError]")
	// take the err and create a new error body
	errBody := errorBody{
		Error: err.Error(), 
	}
	// call respondWithJSON
	respondWithJSON(w, code, errBody)

}

type apiConfig struct {
	dbClient database.Client
}

func (apiCfg apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	// get email via path parameter
	userEmail := strings.TrimPrefix(r.URL.Path, "/users/")
	
	// get other params via path body	
	type parameters struct {
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	// unmarshal parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// call db to update user
	user, err := apiCfg.dbClient.UpdateUser(userEmail, params.Password, params.Name, params.Age)
	if err != nil {
		// update failed
		respondWithError(w, http.StatusBadRequest, err)
		return 
	}

	// update success
	// -> respond with marshalled updated user body
	dat, _ := json.Marshal(user)
	respondWithJSON(w, http.StatusOK, dat)
}



func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	// get email from path 
	userEmail := strings.TrimPrefix(r.URL.Path, "/users/")
	// db call to get user
	user, err := apiCfg.dbClient.GetUser(userEmail)
	
	if err != nil {
		// get failed
		respondWithError(w, http.StatusBadRequest, err)
		return 
	}

	// get success 
	// -> respond with marshalled user body
	dat, _ := json.Marshal(user)
	respondWithJSON(w, http.StatusOK, dat)
}


func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	// get email from path 
	userEmail := strings.TrimPrefix(r.URL.Path, "/users/")
	// db call to delete
	err := apiCfg.dbClient.DeleteUser(userEmail)
	if err != nil {
		// delete failed
		respondWithError(w, http.StatusBadRequest, err)
		return 

	}
	// delete success 
	// -> respond with empty json body
	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (apiCfg apiConfig) handlerCreateUser(
	w http.ResponseWriter, 
	r *http.Request) {
	
	log.Println("IN[handlerCreateUser]")

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	// unmarshal parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	// call db to create new user
	user, err := apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return 
	}
	// db create success
	// respond with new user json
	dat, _ := json.Marshal(user)
	respondWithJSON(w, http.StatusCreated, dat)
}

func (apiCfg apiConfig) endpointUsersHandler(
	w http.ResponseWriter, r *http.Request) {

	log.Println("IN[endpointUsersHandler]")
	log.Println("r.Method: ", r.Method)
	
	switch r.Method {
	case http.MethodGet:
		// call GET handler
		apiCfg.handlerGetUser(w, r)
		break

	case http.MethodPost:
		// call POST handler
		apiCfg.handlerCreateUser(w, r)
		break

	case http.MethodPut:
		// call PUT handler
		apiCfg.handlerUpdateUser(w, r)
		break
	
	case http.MethodDelete:
		// call DELETE handler
		apiCfg.handlerDeleteUser(w, r)
		break

	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}


func (apiCfg apiConfig) endpointPostsHandler(
	w http.ResponseWriter, r *http.Request) {

	log.Println("IN[endpointPostsHandler]")
	log.Println("r.Method: ", r.Method)
	
	switch r.Method {
	case http.MethodGet:
		// call GET handler
		// apiCfg.handlerGetUser(w, r)
		break

	case http.MethodPost:
		// call POST handler
		apiCfg.handlerCreatePost(w, r)
		break

	case http.MethodPut:
		// call PUT handler
		// apiCfg.handlerUpdateUser(w, r)
		break
	
	case http.MethodDelete:
		// call DELETE handler
		apiCfg.handlerDeletePost(w, r)
		break

	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}
func (apiCfg apiConfig) handlerRetrievePosts(w http.ResponseWriter, r *http.Request) {
	
	// get all posts for a specified user
	type parameters struct {
		UserEmail string `json:"userEmail"`
	}

	// unmarshal parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// call db to get all posts for specified user
	userPosts, err := apiCfg.dbClient.GetPosts(params.UserEmail)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return 
	}

	// marshal posts into json
	dat, _ := json.Marshal(userPosts)
	respondWithJSON(w, http.StatusOK, dat)
}

func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
		// get post id from path 
		postId := strings.TrimPrefix(r.URL.Path, "/posts/")
		// db call to delete
		err := apiCfg.dbClient.DeletePost(postId)
		if err != nil {
			// delete failed
			respondWithError(w, http.StatusBadRequest, err)
			return 
	
		}
		// delete success 
		// -> respond with empty json body
		respondWithJSON(w, http.StatusOK, struct{}{})
}

func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserEmail string `json:"userEmail"`
		Text      string `json:"text"`
	}

	// unmarshal parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	// call db to create new post
	post, err := apiCfg.dbClient.CreatePost(params.UserEmail, params.Text)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return 
	}
	// db create success
	// respond with new post json
	dat, _ := json.Marshal(post)
	respondWithJSON(w, http.StatusCreated, dat)	
}

func main() {

	// create a new db client
	dbClient := database.NewClient("./db.json")

	// create an api config 
	dbApiConfig := apiConfig{dbClient: dbClient}

	// create mux
	serveMux := http.NewServeMux()

	// create handler(s)
	var patternRoot string = "/"
	serveMux.HandleFunc(patternRoot, testHandler)
	
	serveMux.HandleFunc("/err", testErrHandler)
	
	// register endpoints-user handler(s)
	serveMux.HandleFunc("/users", dbApiConfig.endpointUsersHandler)
	serveMux.HandleFunc("/users/", dbApiConfig.endpointUsersHandler)

	// register endpoints-post handler(s)
	serveMux.HandleFunc("/posts", dbApiConfig.endpointPostsHandler)
	serveMux.HandleFunc("/posts/", dbApiConfig.endpointPostsHandler)

	// create server
	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      serveMux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	// call server
	srv.ListenAndServe()
}



