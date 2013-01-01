package main

import (
	"fmt"
	"net/rpc"
)

type GitServer struct{}

type Requst struct {
	Path string
	Args []string
}

func (s *GitServer) Result(request *Requst, reply *string) error {
	*reply = ""
	repo := DB.GetRepo(request.Path)
	if repo == nil {
		return fmt.Errorf("can't find repo\n")
	}
	*reply = repo.Result(request.Args)
	return nil
}

func NewServer() {
	server := &GitServer{}
	rpc.Register(server)
	rpc.HandleHTTP()
}
