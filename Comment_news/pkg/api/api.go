package api

import (
	"Comment_news/pkg/dbcomment"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var mat []string = []string{"qwerty", "йцукен", "zxvbnm"}

type API struct {
	db *dbcomment.DB
	r  *mux.Router
}

// Конструктор API.
func New(db *dbcomment.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a
}

// Router HTTP-сервера
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/comments/all", api.comments).Methods(http.MethodGet, http.MethodOptions).
		Queries("newsId", "{newsId:[0-9]+}")
	api.r.HandleFunc("/comments/add", api.addComment).Methods(http.MethodPost, http.MethodOptions)

}

func (api *API) comments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	fmt.Println("ghere")
	s := mux.Vars(r)["newsId"]
	id, _ := strconv.Atoi(s)
	comments, err := api.db.Comment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(comments); i++ {
		if comments[i].Profane == true {
			comments[i] = comments[len(comments)-1]
			comments = comments[:len(comments)-1]
		}

	}
	json.NewEncoder(w).Encode(comments)
}

func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var comment dbcomment.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(comment, err)

	for i := 0; i < len(mat); i++ {
		if profane := strings.Contains(comment.Text, mat[i]); profane == true {
			comment.Profane = true
		}

	}
	err = api.db.AddComment(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
