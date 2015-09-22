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
	RoomId      string          `json:"room_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Messages    []MessageObject `json:"messages"`
}

type RoomSetResource struct {
	db            *sql.DB
	getStatement  *sql.Stmt
	postStatement *sql.Stmt
}

type RoomResource struct {
	db                  *sql.DB
	getStatement        *sql.Stmt
	getMessageStatement *sql.Stmt
}

func NewRoomSetResource(db *sql.DB) *RoomSetResource {
	getStatement, err := db.Prepare(
		"SELECT room_id, name, description FROM room;",
	)
	fatal(err)
	postStatement, err := db.Prepare(
		"SELECT success, room_id FROM add_room($1, $2);", // name, description
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
	getMessageStatement, err := db.Prepare(
		"SELECT message_id, account_id, room_id, name, body, " +
			"to_char(message.created, 'Dy Mon DD HH24:MI:SS TZ YYYY') FROM " +
			"message JOIN account USING(account_id) " +
			"WHERE room_id=$1 ORDER BY message_id;",
	)
	fatal(err)
	return &RoomResource{db, getStatement, getMessageStatement}
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

	// Room
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

	// Room Messages
	rows2, err := r.getMessageStatement.Query(result.RoomId)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows2.Close()
	for rows2.Next() {
		var messageObject MessageObject
		err = rows2.Scan(&messageObject.MessageId, &messageObject.AccountId, &messageObject.RoomId, &messageObject.Name, &messageObject.Body, &messageObject.Created)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else {
			result.Messages = append(result.Messages, messageObject)
		}
	}
	err = rows2.Err()
	if err != nil {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	return http.StatusOK, result, nil
}
