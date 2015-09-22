package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	gogoapi "github.com/Sennue/gogoapi"
)

type RoomObject struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomsResource struct {
	db           *sql.DB
	getStatement *sql.Stmt
}

type RoomResource struct {
	db           *sql.DB
	getStatement *sql.Stmt
}

func NewRoomsResource(db *sql.DB) *RoomsResource {
	getStatement, err := db.Prepare(
		"SELECT room_id, name, description FROM room;",
	)
	fatal(err)
	return &RoomsResource{db, getStatement}
}

func (r *RoomsResource) Get(request *http.Request) (int, interface{}, http.Header) {
	var result []RoomObject

	rows, err := r.getStatement.Query()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows.Close()
	for rows.Next() {
		var (
			room_id     int
			name        string
			description string
		)
		err = rows.Scan(&room_id, &name, &description)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else {
			result = append(result, RoomObject{room_id, name, description})
		}
	}
	err = rows.Err()
	if err != nil {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	return http.StatusOK, result, nil
}

func NewRoomResource(db *sql.DB) *RoomResource {
	getStatement, err := db.Prepare(
		"SELECT room_id, name, description FROM room " +
			"WHERE room_id=$1",
	)
	fatal(err)
	return &RoomResource{db, getStatement}
}

func (r *RoomResource) Get(request *http.Request) (int, interface{}, http.Header) {
	var (
		result RoomObject
		success bool = false
	)

	vars := mux.Vars(request)
	room_id := vars["room_id"]

	rows, err := r.getStatement.Query(
		room_id,
	)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows.Close()
	for rows.Next() {
		var (
			room_id     int
			name        string
			description string
		)
		err = rows.Scan(&room_id, &name, &description)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else {
			result = RoomObject{room_id, name, description}
			success = true
		}
	}
	err = rows.Err()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	if !success {
		// 204 StatusNoContent does not return a body
		status := http.StatusNoContent
		message := fmt.Sprintf("No room with id %s.", room_id)
		return status, gogoapi.JSONMessage{status, message}, nil
	}
	return http.StatusOK, result, nil
}
