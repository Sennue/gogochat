package main

import (
	"log"
)

func fatal(err error) {
	if nil != err {
		log.Fatal(err)
	}
}
