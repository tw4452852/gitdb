package main

import (
	"container/list"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

var (
	DB   *repoDb
	help *bool
	root *string
)

const (
	gitdir = ".git"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//flags
	help = flag.Bool("h", false, "show help")
	root = flag.String("root", "/", "set the root directory for searching")
}

type repoDb struct {
	repos *list.List
}

func (db *repoDb) GetRepo(path string) *gitRepo {
	for e := db.repos.Front(); e != nil; e = e.Next() {
		if r := e.Value.(*gitRepo); r.match(path) {
			return r
		}
	}
	return nil
}

func (db *repoDb) loop(pathC <-chan string) {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case path := <-pathC:
			db.repos.PushBack(NewRepo(path))
		case <-timer.C:
			db.update()
		}
	}
}

func (db *repoDb) update() {
	cleanup := []*list.Element{}
	for e := db.repos.Front(); e != nil; e = e.Next() {
		if r := e.Value.(*gitRepo); r.exist() {
			r.ExitC <- struct{}{}
			cleanup = append(cleanup, e)
		}
	}
	for _, e := range cleanup {
		db.repos.Remove(e)
	}
}

func NewRepoDb(pathC <-chan string) *repoDb {
	db := &repoDb{
		repos: list.New(),
	}
	go db.loop(pathC)
	return db
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	DB = NewRepoDb(NewSpider(*root))
	InitServer()

	err := http.ListenAndServe(":54321", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
