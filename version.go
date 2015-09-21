package main

import (
	"net/http"
)

type VersionObject struct {
	StatusCode int    `json:"status"`
	Version    string `json:"version"`
}

func (v VersionObject) Status() int {
	return v.StatusCode
}

type VersionResource struct {
	Version string
}

func NewVersionResource(version string) *VersionResource {
	return &VersionResource{version}
}

func (v *VersionResource) Get(request *http.Request) (int, interface{}, http.Header) {
	status := http.StatusOK
	return status, VersionObject{status, v.Version}, nil
}
