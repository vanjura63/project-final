package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var wg sync.WaitGroup = sync.WaitGroup{}

type Comment struct {
	Id       int    `json:"id"`
	NewsId   int    `json:"newsId"`
	Text     string `json:"text"`
	ParentId int    `json:"parentId"`
	PubTime  int    `json:"pubTime"`
	Profane  bool   `json:"-"`
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pubTime"`
	Link    string `json:"link"`
}

type PostDetailed struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	PubTime int64     `json:"pubTime"`
	Link    string    `json:"link"`
	Comment []Comment `json:"comments"`
}

// Application configuration
type Config struct {
	Servnews    string `json:"servnews"`
	Servcomment string `json:"servcomment"`
}

type API struct {
	r      *mux.Router
	config Config
}

var commentsChan chan []Comment = make(chan []Comment, 500)

// Конструктор API.
func New(c Config) *API {
	a := API{r: mux.NewRouter()}
	a.endpoints()
	a.config = c
	return &a
}

// Router HTTP-сервера
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {

	api.r.HandleFunc("/news", api.posts).Methods(http.MethodGet, http.MethodOptions).
		Queries("n", "{n:[0-9]+}")
	api.r.HandleFunc("/news/filter", api.filterNews).Methods(http.MethodGet, http.MethodOptions).
		Queries("query", "{query:.+}")
	/*api.r.HandleFunc("/news/comments", api.posts).Methods(http.MethodGet, http.MethodOptions).
	Queries("newsId", "{newsId:[0-9]+}")*/
	api.r.HandleFunc("/news/comments/add", api.addComment).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/news/single", api.newsSingleDetailed).Methods(http.MethodGet, http.MethodOptions).
		Queries("newsId", "{newsId:[0-9]+}")

}

// Получение новостного блока
func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	news, err := http.Get(api.config.Servnews + "/news?n=" + s)

	fmt.Println("Новостиии--", news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var newsdecod []PostDetailed
	err = json.NewDecoder(news.Body).Decode(&newsdecod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wg.Add(len(newsdecod))
	for i := 0; i < len(newsdecod); i++ {
		go func(i int) {
			defer wg.Done()
			if newsdecod[i].ID > 0 {
				fmt.Println(api.config.Servcomment + "/comments/all?newsId=" + fmt.Sprintf("%d", newsdecod[i].ID))
				resp, err := http.Get(api.config.Servcomment + "/comments/all?newsId=" + fmt.Sprintf("%d", newsdecod[i].ID))
				fmt.Println("Kомент:", err)
				var comments []Comment
				err = json.NewDecoder(resp.Body).Decode(&comments)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				commentsChan <- comments

			}

		}(i)

	}
	wg.Wait()
	close(commentsChan)
	for c := range commentsChan {
		if c != nil && len(c) > 0 {
			for i := 0; i < len(newsdecod); i++ {
				if newsdecod[i].ID == c[0].NewsId {
					newsdecod[i].Comment = c
				}
			}
		}
	}
	commentsChan = make(chan []Comment, 500)
	json.NewEncoder(w).Encode(newsdecod)

}

// Поиск конкретной новости
func (api *API) filterNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var fnews []Post
	s := mux.Vars(r)["query"]
	if s == "" {
		http.Error(w, "Query не может быть пустым", http.StatusInternalServerError)
		return
	}
	news, err := http.Get(api.config.Servnews + "/news/filter?query=" + s)

	err = json.NewDecoder(news.Body).Decode(&fnews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(fnews)
}

// Добавление коментария к новости
func (api *API) addComment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(comment)
	resp, err := http.Post(api.config.Servcomment+"/comments/add", "application/json", bytes.NewBuffer(body))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)

}

// Определение детальной новости
func (api *API) newsSingleDetailed(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	var postd PostDetailed
	s := mux.Vars(r)["newsId"]
	//id, _ := strconv.Atoi(s)
	idnews, err := http.Get(api.config.Servnews + "/news/single?idnews=" + s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("ID news:", idnews)

	err = json.NewDecoder(idnews.Body).Decode(&postd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var comments []Comment

	fmt.Println("Comment:", comments)

	if postd.ID > 0 {
		fmt.Println(api.config.Servcomment + "/comments/all?newsId" + fmt.Sprintf("%d", postd.ID))
		resp, err := http.Get(api.config.Servcomment + "/comments/all?newsId=" + fmt.Sprintf("%d", postd.ID))
		fmt.Println("Kомент:", resp)
		err = json.NewDecoder(resp.Body).Decode(&comments)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	postd.Comment = comments

	json.NewEncoder(w).Encode(postd)

}
