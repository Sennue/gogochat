package main

import (
	"net/http"

	gogoapi "github.com/Sennue/gogoapi"
)

type PingResource struct {
}

func (ping *PingResource) Pong() (int, interface{}, http.Header) {
	status := http.StatusOK
	return status, gogoapi.JSONMessage{status, "pong"}, nil
}

func NewPingResource() *PingResource {
	return &PingResource{}
}

func (ping *PingResource) Get(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}

func (ping *PingResource) Post(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}

func (ping *PingResource) Put(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}

func (ping *PingResource) Delete(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}

func (ping *PingResource) Head(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}

func (ping *PingResource) Patch(request *http.Request) (int, interface{}, http.Header) {
	return ping.Pong()
}
