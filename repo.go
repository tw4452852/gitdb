package main

import (
	"container/list"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type cmdEntry struct {
	args   []string
	result string
}

func (e *cmdEntry) update(root string) error {
	cmd := exec.Command("git", e.args...)
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	e.result = string(out)
	return nil
}

//don't care about the arg sequence
func (e *cmdEntry) match(args []string) bool {
	if len(args) != len(e.args) {
		return false
	}
FINDING:
	for _, dst := range args {
		for _, src := range e.args {
			if dst == src {
				//found
				continue FINDING
			}
		}
		//not found
		return false
	}
	return true
}

type gitRepo struct {
	ExitC   chan struct{}
	root    string
	addC    chan *cmdEntry
	reqC    chan *request
	entries *list.List
}

type request struct {
	args  []string
	reply chan string
}

//get someone result, if not exist, add it
func (r *gitRepo) Result(args []string) string {
	req := &request{
		args:  args,
		reply: make(chan string),
	}
	r.reqC <- req
	return <-req.reply
}

func (r *gitRepo) match(path string) bool {
	return strings.Contains(path, r.root)
}

func (r *gitRepo) exist() bool {
	_, err := os.Lstat(r.root + gitdir)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

//used for update/add/remove every cmd entry
func (r *gitRepo) loop() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case e := <-r.addC:
			r.entries.PushBack(e)
		case req := <-r.reqC:
			found := false
			for e := r.entries.Front(); e != nil; e = e.Next() {
				if e.Value.(*cmdEntry).match(req.args) {
					req.reply <- e.Value.(*cmdEntry).result
					found = true
					break
				}
			}
			if !found {
				r.addC <- &cmdEntry{req.args, ""}
				req.reply <- ""
			}
		case <-timer.C:
			r.update()
		case <-r.ExitC:
			//this repo isn't exist
			return
		}
	}
}

func (r *gitRepo) update() {
	cleanup := []*list.Element{}
	waiter := sync.WaitGroup{}
	for e := r.entries.Front(); e != nil; e = e.Next() {
		e := e
		go func() {
			waiter.Add(1)
			if err := e.Value.(*cmdEntry).update(r.root); err != nil {
				cleanup = append(cleanup, e)
			}
			waiter.Done()
		}()
	}
	waiter.Wait()
	for _, e := range cleanup {
		r.entries.Remove(e)
	}
}

func NewRepo(root string) *gitRepo {
	repo := &gitRepo{
		ExitC:   make(chan struct{}),
		root:    root,
		addC:    make(chan *cmdEntry, 10),
		reqC:    make(chan *request),
		entries: list.New(),
	}
	log.Println("NewRepo:", root)
	go repo.loop()
	return repo
}
