package main

import (
	"fmt"
)

func main() {
	pathC := NewSpider("") //frome root directory

	for path := range pathC {
		fmt.Printf("get a repo, path is %q\n", path)
	}
}
