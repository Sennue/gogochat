package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	gogoapi "github.com/Sennue/gogoapi"
)

type MessageObject struct {
	MessageId string `json:"message_id"`
	AccountId string `json:"account_id"`
	RoomId    string `json:"room_id"`
	Name      string `json:"name"`
	Body      string `json:"body"`
	Created   string `json:"created"`
}

type MessageSetResource struct {
	auth          *gogoapi.AuthResource
	db            *sql.DB
	postStatement *sql.Stmt
}

type MessageResource struct {
	db           *sql.DB
	getStatement *sql.Stmt
}

func NewMessageSetResource(auth *gogoapi.AuthResource, db *sql.DB) *MessageSetResource {
	postStatement, err := db.Prepare(
		"SELECT success, message_id FROM add_message($1, $2, $3);", // account_id, room_id, body
	)
	fatal(err)
	return &MessageSetResource{auth, db, postStatement}
}

func (r *MessageSetResource) Post(request *http.Request) (int, interface{}, http.Header) {
	var (
		messageObject MessageObject
		success       bool
	)

	authorized, token, err := r.auth.IsAuthorized(request)
	if !authorized {
		return gogoapi.AuthorizationFailed()
	}

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	err = request.Body.Close()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	if err := json.Unmarshal(body, &messageObject); nil != err {
		return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
	}
	_, err = strconv.ParseInt(messageObject.RoomId, 10, 64)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	switch customUserInfo := token.Claims["CustomUserInfo"].(type) {
	case map[string]interface{}:
		switch accountId := customUserInfo["account_id"].(type) {
		case float64:
			messageObject.AccountId = strconv.FormatInt(int64(accountId), 10)
		default:
			log.Printf("account_id type %T %t", accountId, accountId)
			return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
		}
		switch username := customUserInfo["username"].(type) {
		case string:
			messageObject.Name = username
		default:
			log.Printf("username type %T %t", username, username)
			return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
		}
	default:
		log.Printf("customUserInfo type %T %t", customUserInfo, customUserInfo)
		return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
	}

	err = r.postStatement.QueryRow(
		messageObject.AccountId,
		messageObject.RoomId,
		messageObject.Body,
	).Scan(&success, &messageObject.MessageId)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	if !success {
		status := http.StatusConflict
		return status, gogoapi.JSONError{status, "Resource already exists."}, nil
	} else {
		messageObject.Created = fmt.Sprintf(time.Now().Format(time.UnixDate))
		return http.StatusCreated, messageObject, nil
	}
}

func NewMessageResource(db *sql.DB) *MessageResource {
	getStatement, err := db.Prepare(
		"SELECT message_id, account_id, room_id, name, body, " +
			"to_char(message.created, 'Dy Mon DD HH24:MI:SS TZ YYYY') FROM " +
			"message JOIN account USING(account_id) WHERE message_id=$1;",
	)
	fatal(err)
	return &MessageResource{db, getStatement}
}

func (r *MessageResource) Get(request *http.Request) (int, interface{}, http.Header) {
	var (
		result  MessageObject
		success bool = false
	)

	vars := mux.Vars(request)
	message_id := vars["message_id"]

	_, err := strconv.ParseInt(message_id, 10, 64)
	if nil != err {
		return http.StatusNoContent, nil, nil
	}

	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	rows, err := r.getStatement.Query(
		message_id,
	)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&result.MessageId, &result.AccountId, &result.RoomId, &result.Name, &result.Body, &result.Created)
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
		message := fmt.Sprintf("No message with id %s.", message_id)
		return status, gogoapi.JSONMessage{status, message}, nil
	}
	return http.StatusOK, result, nil
}
