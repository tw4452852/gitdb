package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Requst struct {
	Path string
	Args []string
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: ", os.Args[0], "something")
		os.Exit(1)
	}
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:54321")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("getwd:", err)
	}
	request := &Requst{
		Path: pwd,
		Args: os.Args[1:],
	}
	var reply string
	err = client.Call("GitServer.Result", request, &reply)
	if err != nil {
		log.Fatal("call:", err)
	}
	fmt.Print(reply)
}
