package main

import (
	"fmt"
	"net/http"
	"time"
)

type TimeObject struct {
	StatusCode int    `json:"status"`
	Time       string `json:"time"`
}

func (t TimeObject) Status() int {
	return t.StatusCode
}

type TimeResource struct {
}

func NewTimeResource() *TimeResource {
	return &TimeResource{}
}

func (t *TimeResource) Get(request *http.Request) (int, interface{}, http.Header) {
	status := http.StatusOK
	time := fmt.Sprintf(time.Now().Format(time.UnixDate))
	return status, TimeObject{status, time}, nil
}
