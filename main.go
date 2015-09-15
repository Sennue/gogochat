package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	gogoapi "github.com/Sennue/gogoapi"
	_ "github.com/lib/pq"
)

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
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

	dbdriver := "postgres"
	dbhost := "localhost"
	dbname := "gogochat"
	dbuser := "gogochat"
	dbpassword := "gogochat"
	dboptions := "sslmode=disable"
	dbDataSource := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s %s",
		dbuser, dbpassword, dbname, dbhost, dboptions,
	)
	db, err := sql.Open(dbdriver, dbDataSource)
	fatal(err)
	defer db.Close()

	api := gogoapi.NewAPI([]gogoapi.WrapperFunc{gogoapi.Logger})
	user := NewUserResource(1, "user", "password", "")
	validator := NewAuthValidator(db)
	auth := gogoapi.NewAuthResource(privateKeyPath, publicKeyPath, 60, validator.Validate)
	api.AddResource(auth, "/auth", nil)
	api.AddResource(user, "/user", []gogoapi.WrapperFunc{auth.AuthorizationRequired})
	if err := api.Start(host, port); nil != err {
		log.Fatal(err)
	}
}
