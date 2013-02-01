package main

import (
	"log"
	"net/http"
)

var DB *repoDb

func main() {
	DB = NewRepoDb(NewSpider("")) //frome root directory

	err := http.ListenAndServe("127.0.0.1:54321", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
