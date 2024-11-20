package main

import (
	"flag"
	"fmt"
	"runtime"

	camect_go "github.com/kaiiorg/camect-go"
)

var (
	ip       = flag.String("ip", "0.0.0.0", "ip address of hub")
	username = flag.String("username", "admin", "username of hub local admin")
	password = flag.String("password", "this isn't a real password, provide your own", "password of hub local admin")
)

func main() {
	flag.Parse()

	fmt.Printf("camect-go on %s\n", runtime.Version())

	hub := camect_go.New(*ip, *username, *password)

	info, err := hub.Info()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hub Name: %s\n", info.Name)
	fmt.Printf("Hub ID: %s\n", info.Id)
}
