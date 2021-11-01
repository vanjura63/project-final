package main

import (
	"Go_news/pkg/api"
	"Go_news/pkg/dbnews"
	"Go_news/pkg/rss"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Application configuration
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {

	connstr := "postgres://postgres:ts950sdx@localhost/dbnews"

	// Инициализация БД

	db, err := dbnews.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// Чтение файла конфигурации

	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Файл конфига:", string(b))

	//Разбор новостей в потоке

	chPosts := make(chan []dbnews.Post)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	// Запись новостей потока в БД
	go func() {
		for posts := range chPosts {
			db.StoreNews(posts)
		}
	}()

	// Контроль ошибок
	go func() {
		for err := range chErrs {
			log.Println("Errors:", err)
		}
	}()

	//Старт WEB сервера
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

// Поток каналов
func parseURL(url string, db *dbnews.DB, posts chan<- []dbnews.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
