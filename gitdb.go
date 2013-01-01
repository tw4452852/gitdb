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
	DB      *repoDb
	rootDir *string
	help    *bool
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//flags
	rootDir = flag.String("root", "/", "base directory where spider begin with")
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
FINDROOT:
	for {
		if dir == "." {
			log.Printf("addRepo: can't find repo root directory with path(%q)\n", path)
			return
		}
		f, err := os.Open(dir)
		if err != nil {
			log.Println("addRepo:", err)
			return
		}
		names, err := f.Readdirnames(-1)
		if err != nil {
			log.Println("addRepo:", err)
			return
		}
		for _, name := range names {
			if name == gitdir {
				break FINDROOT
			}
		}
		f.Close()
		dir = filepath.Dir(dir)
	}
	db.repos = append(db.repos, NewRepo(dir))
}

func NewRepoDb(pathC <-chan string) *repoDb {
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

	DB = NewRepoDb(NewSpider(*rootDir)) //frome root directory
	NewServer()

	err := http.ListenAndServe(":54321", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
