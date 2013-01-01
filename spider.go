package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type spiderDb struct {
	base   string
	output chan string
	log    map[string]struct{}
}

func NewSpider(base string) <-chan string {
	if base == "" {
		base = "/"
	}
	db := &spiderDb{
		base:   base,
		output: make(chan string),
		log:    make(map[string]struct{}),
	}
	go spiderLoop(db)
	return db.output
}

func spiderLoop(db *spiderDb) {
	for {
		err := filepath.Walk(db.base,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return filepath.SkipDir
				}
				if !info.IsDir() {
					return nil
				}
				if info.Name() == ".git" {
					addGitRepo(filepath.Dir(path), db)
					return filepath.SkipDir
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func addGitRepo(path string, db *spiderDb) {
	if _, found := db.log[path]; found {
		return
	}

	db.log[path] = struct{}{}
	db.output <- path
}
