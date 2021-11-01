package api

import (
	"Go_news/pkg/dbnews"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	db *dbnews.DB
	r  *mux.Router
}

// Конструктор API.
func New(db *dbnews.DB) *API {
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
	api.r.HandleFunc("/news", api.posts).Methods(http.MethodGet, http.MethodOptions).Queries("n", "{n:[0-9]+}")
	// Получить одну новость
	api.r.HandleFunc("/news/single", api.article).Methods(http.MethodGet, http.MethodOptions).Queries("idnews", "{idnews:[0-9]+}")
	api.r.HandleFunc("/news/filter", api.newsfilter).Methods(http.MethodGet, http.MethodOptions).Queries("query", "{query:.+}")

}

func (api *API) newsfilter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["query"]
	fmt.Println("Query---", s)
	news, err := api.db.FilterNews(s)
	fmt.Println("Получено из БД", err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	news, err := api.db.News(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}

func (api *API) article(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	s := mux.Vars(r)["idnews"]
	fmt.Println("IDnews:", s)
	id, _ := strconv.Atoi(s)

	fmt.Println("ID:", id)

	news, err := api.db.Article(id)
	fmt.Println(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}
