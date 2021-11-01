package main

import (
	"Comment_news/pkg/api"
	"Comment_news/pkg/dbcomment"
	"log"
	"net/http"
)

func main() {

	connstr := "postgres://postgres:ts950sdx@localhost/dbcomment"

	// Init DB

	db, err := dbcomment.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// Web server start
	err = http.ListenAndServe(":8181", api.Router())
	if err != nil {
		log.Fatal(err)
	}

}
