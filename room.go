package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	gogoapi "github.com/Sennue/gogoapi"
)

type RoomObject struct {
	RoomId      string `json:"room_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomSetResource struct {
	db            *sql.DB
	getStatement  *sql.Stmt
	postStatement *sql.Stmt
}

type RoomResource struct {
	db           *sql.DB
	getStatement *sql.Stmt
}

func NewRoomSetResource(db *sql.DB) *RoomSetResource {
	getStatement, err := db.Prepare(
		"SELECT room_id, name, description FROM room;",
	)
	fatal(err)
	postStatement, err := db.Prepare(
		"SELECT success, account_id FROM add_room($1, $2);",
		//"INSERT INTO room (name, description) VALUES ($1, $2);",
	)
	fatal(err)
	return &RoomSetResource{db, getStatement, postStatement}
}

func (r *RoomSetResource) Get(request *http.Request) (int, interface{}, http.Header) {
	var result []RoomObject

	rows, err := r.getStatement.Query()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows.Close()
	for rows.Next() {
		var roomObject RoomObject
		err = rows.Scan(&roomObject.RoomId, &roomObject.Name, &roomObject.Description)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else {
			result = append(result, roomObject)
		}
	}
	err = rows.Err()
	if err != nil {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	return http.StatusOK, result, nil
}

func (r *RoomSetResource) Post(request *http.Request) (int, interface{}, http.Header) {
	var (
		roomObject RoomObject
		success    bool
	)

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	err = request.Body.Close()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	if err := json.Unmarshal(body, &roomObject); nil != err {
		return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
	}

	err = r.postStatement.QueryRow(
		roomObject.Name,
		roomObject.Description,
	).Scan(&success, &roomObject.RoomId)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	if !success {
		status := http.StatusConflict
		return status, gogoapi.JSONError{status, "Resource already exists."}, nil
	} else {
		return http.StatusCreated, roomObject, nil
	}
}

func NewRoomResource(db *sql.DB) *RoomResource {
	getStatement, err := db.Prepare(
		"SELECT room_id, name, description FROM room " +
			"WHERE room_id=$1;",
	)
	fatal(err)
	return &RoomResource{db, getStatement}
}

func (r *RoomResource) Get(request *http.Request) (int, interface{}, http.Header) {
	var (
		result  RoomObject
		success bool = false
	)

	vars := mux.Vars(request)
	room_id := vars["room_id"]

	_, err := strconv.ParseInt(room_id, 10, 64)
	if nil != err {
		return http.StatusNoContent, nil, nil
	}

	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	rows, err := r.getStatement.Query(
		room_id,
	)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.RoomId, &result.Name, &result.Description)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else {
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
