package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: ", os.Args[0], "st/status")
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
	var reply string
	err = client.Call("GitServer.Status", pwd, &reply)
	if err != nil {
		log.Fatal("call:", err)
	}
	fmt.Print(reply)
}
