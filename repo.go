package main

import (
	"log"
	"os/exec"
	"strings"
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
	root    string
	addC    chan *cmdEntry
	reqC    chan *request
	entries []*cmdEntry
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

//used for update/add/remove every cmd entry
func (r *gitRepo) loop() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case e := <-r.addC:
			r.entries = append(r.entries, e)
		case req := <-r.reqC:
			found := false
			for _, e := range r.entries {
				if e.match(req.args) {
					req.reply <- e.result
					found = true
					break
				}
			}
			if !found {
				r.addC <- &cmdEntry{req.args, ""}
				req.reply <- ""
			}
		case <-timer.C:
			cleanup := []int{}
			for i, entry := range r.entries {
				if err := entry.update(r.root); err != nil {
					cleanup = append(cleanup, i)
				}
			}
			for _, i := range cleanup {
				r.entries = append(r.entries[:i], r.entries[i+1:]...)
			}
		}
	}
}

func NewRepo(root string) *gitRepo {
	repo := &gitRepo{
		root:    root,
		addC:    make(chan *cmdEntry, 10),
		reqC:    make(chan *request),
		entries: make([]*cmdEntry, 0, 10),
	}
	log.Println("NewRepo:", root)
	go repo.loop()
	return repo
}
