package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

type gitRepo struct {
	root   string
	status string
}

func (r *gitRepo) Status() string {
	return r.status
}

func (r *gitRepo) match(path string) bool {
	return strings.Contains(path, r.root)
}

func (r *gitRepo) loop() {
	timer := time.NewTicker(1 * time.Second)
	for _ = range timer.C {
		r.getStatus()
	}
}

func (r *gitRepo) getStatus() {
	cmd := exec.Command("git", "status")
	//set the working path
	cmd.Dir = r.root
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return
	}
	//update status
	r.status = string(out)
}

func NewRepo(root string) *gitRepo {
	repo := &gitRepo{
		root: root,
	}
	go repo.loop()
	return repo
}

type repoDb struct {
	repos []*gitRepo
}

func (db *repoDb) getRepo(path string) *gitRepo {
	for _, r := range db.repos {
		if r.match(path) {
			return r
		}
	}
	return nil
}
func (db *repoDb) loop(pathC <-chan string) {
	for path := range pathC {
		r := NewRepo(path)
		db.repos = append(db.repos, r)
	}
}

func NewRepoDb(pathC <-chan string) *repoDb {
	db := &repoDb{
		repos: make([]*gitRepo, 0, 10),
	}
	go db.loop(pathC)
	return db
}
