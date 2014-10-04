package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kelseyhightower/pm/cmd/build"
	"github.com/kelseyhightower/pm/cmd/get"
	"github.com/kelseyhightower/pm/cmd/info"
	"github.com/kelseyhightower/pm/cmd/install"
	"github.com/kelseyhightower/pm/cmd/list"
	"github.com/kelseyhightower/pm/cmd/remove"
	"github.com/kelseyhightower/pm/cmd/verify"
)

func main() {
	log.SetFlags(0)
	switch os.Args[1] {
	case "build":
		build.Run()
	case "verify":
		verify.Run()
	case "install":
		install.Run()
	case "get":
		get.Run()
	case "remove":
		remove.Run()
	case "list":
		list.Run()
	case "info":
		info.Run()
	default:
		fmt.Println("not a command")
	}
}
