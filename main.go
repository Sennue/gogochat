package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"

	gogoapi "github.com/Sennue/gogoapi"
	_ "github.com/lib/pq"
)

const (
	API_VERSION          = "1.0.0"
	DEFAULT_PROGRAM_NAME = "gogochat"
	DEFAULT_HOST         = "localhost"
	DEFAULT_PORT         = 8080
	DEFAULT_PUBLIC_KEY   = "keys/app.rsa.pub"
	DEFAULT_PRIVATE_KEY  = "keys/app.rsa"
)

func initSyslog(tag string) {
	logwriter, err := syslog.New(syslog.LOG_NOTICE, tag)
	fatal(err)
	log.SetOutput(logwriter)
}

func main() {
	var (
		host                               string
		port                               int
		syslog                             bool
		syslogTag                          string
		dbhost, dbname, dbuser, dbpassword string
		publicKey, privateKey              string
	)
	flag.StringVar(&host, "host", DEFAULT_HOST, "host address to listen on")
	flag.IntVar(&port, "port", DEFAULT_PORT, "port to listen on")
	flag.BoolVar(&syslog, "syslog", true, "use syslog for logging")
	flag.StringVar(&syslogTag, "tag", DEFAULT_PROGRAM_NAME, "tag used syslog for logging")
	flag.StringVar(&dbhost, "dbhost", DEFAULT_HOST, "database host")
	flag.StringVar(&dbname, "dbname", DEFAULT_PROGRAM_NAME, "database name")
	flag.StringVar(&dbuser, "duser", DEFAULT_PROGRAM_NAME, "database user")
	flag.StringVar(&dbpassword, "dbnpassword", DEFAULT_PROGRAM_NAME, "database user")
	flag.StringVar(&publicKey, "public", DEFAULT_PUBLIC_KEY, "JSON web token public key")
	flag.StringVar(&privateKey, "private", DEFAULT_PRIVATE_KEY, "JSON web token private key")
	flag.Parse()
	if syslog {
		initSyslog(syslogTag)
	}

	// initialize database
	dbdriver := "postgres"
	dboptions := "sslmode=disable"
	dbDataSource := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s %s",
		dbuser, dbpassword, dbname, dbhost, dboptions,
	)
	db, err := sql.Open(dbdriver, dbDataSource)
	fatal(err)

	// cleanup routine
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		db.Close()
		log.Printf("Stopping %s on signal: %v\n", syslogTag, s)
		os.Exit(1)
	}()

	// initialize api
	api := gogoapi.NewAPI([]gogoapi.WrapperFunc{gogoapi.Logger})
	version := NewVersionResource(API_VERSION)
	api.AddResource(version, "/", nil)
	api.AddResource(version, "/version", nil)
	ping := NewPingResource()
	api.AddResource(ping, "/ping", nil)
	time := NewTimeResource()
	api.AddResource(time, "/time", nil)
	validator := NewAuthValidator(db)
	auth := gogoapi.NewAuthResource(privateKey, publicKey, 60, validator.Validate)
	api.AddResource(auth, "/auth", nil)
	account := NewAccountResource(auth, db)
	api.AddResource(account, "/account", nil)
	roomSet := NewRoomSetResource(db)
	api.AddResource(roomSet, "/room", []gogoapi.WrapperFunc{auth.AuthorizationRequired})
	room := NewRoomResource(db)
	api.AddResource(room, "/room/{room_id}", []gogoapi.WrapperFunc{auth.AuthorizationRequired})
	messageSet := NewMessageSetResource(auth, db)
	api.AddResource(messageSet, "/message", nil) // custom auth
	message := NewMessageResource(db)
	api.AddResource(message, "/message/{message_id}", []gogoapi.WrapperFunc{auth.AuthorizationRequired})

	// start server
	log.Printf("Listening on %s:%d\n", host, port)
	if err := api.Start(host, port); nil != err {
		log.Fatal(err)
	}
}
