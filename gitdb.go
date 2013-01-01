package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
)

var (
	DB   *repoDb
	help *bool
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//flags
	help = flag.Bool("h", false, "show help")
}

type repoDb struct {
	repos []*gitRepo
}

func (db *repoDb) GetRepo(path string) *gitRepo {
	for _, r := range db.repos {
		if r.match(path) {
			return r
		}
	}
	db.addRepo(path)
	return nil
}

func (db *repoDb) addRepo(path string) {
	//find the repo root directory
	dir := path
	const gitdir = ".git"
	findGit := func(dir string) bool {
		f, err := os.Open(dir)
		if err != nil {
			log.Println("addRepo:", err)
			return false
		}
		defer f.Close()
		names, err := f.Readdirnames(-1)
		if err != nil {
			log.Println("addRepo:", err)
			return false
		}
		for _, name := range names {
			if name == gitdir {
				return true
			}
		}
		return false
	}

	for {
		if dir == "." || dir == "/" {
			//can't find repo with path
			return
		}
		if findGit(dir) {
			break
		}
		dir = filepath.Dir(dir)
	}
	db.repos = append(db.repos, NewRepo(dir))
}

func NewRepoDb() *repoDb {
	db := &repoDb{
		repos: make([]*gitRepo, 0, 10),
	}
	return db
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	DB = NewRepoDb() //frome root directory
	InitServer()

	err := http.ListenAndServe(":54321", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
