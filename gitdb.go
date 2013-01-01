package main

import (
	"log"
	"net/http"
)

var DB *repoDb

func main() {
	DB = NewRepoDb(NewSpider("")) //frome root directory
	NewServer()

	err := http.ListenAndServe(":54321", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
