package main

import (
	"flag"
	"fmt"
	"log"

	gogoapi "github.com/Sennue/gogoapi"
)

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath = "keys/app.rsa.pub"
)

func main() {
	var (
		host string
		port int
	)
	flag.StringVar(&host, "host", "localhost", "host address to listen on")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()
	fmt.Printf("Listening on %s:%d\n", host, port)

	api := gogoapi.NewAPI([]gogoapi.WrapperFunc{gogoapi.Logger})
	user := NewUserResource(1, "user", "password", "")
	auth := gogoapi.NewAuthResource(privateKeyPath, publicKeyPath, 60)
	api.AddResource(auth, "/auth", nil)
	api.AddResource(user, "/user", []gogoapi.WrapperFunc{auth.AuthorizationRequired})
	if err := api.Start(host, port); nil != err {
		log.Fatal(err)
	}
}
