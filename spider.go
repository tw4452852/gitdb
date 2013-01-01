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
	db := &spiderDb{
		base:   base,
		output: make(chan string),
		log:    make(map[string]struct{}),
	}
	go db.loop()
	return db.output
}

func (db *spiderDb) loop() {
	for {
		err := filepath.Walk(db.base,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return filepath.SkipDir
				}
				if !info.IsDir() {
					return nil
				}
				if info.Name() == gitdir {
					db.addGitRepo(filepath.Dir(path))
					return filepath.SkipDir
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		time.Sleep(30 * time.Second)
	}
}

func (db *spiderDb) addGitRepo(path string) {
	if _, found := db.log[path]; found {
		return
	}

	db.log[path] = struct{}{}
	db.output <- path
}
