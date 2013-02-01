package main

import (
	"fmt"
	"net/rpc"
)

type GitServer struct {
}

func (s *GitServer) Status(path string, reply *string) error {
	*reply = ""
	repo := DB.getRepo(path)
	if repo == nil {
		return fmt.Errorf("can't find repo\n")
	}
	*reply = repo.Status()
	return nil
}

func NewServer() {
	server := &GitServer{}
	rpc.Register(server)
	rpc.HandleHTTP()
}
