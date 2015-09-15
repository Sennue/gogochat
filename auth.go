package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	gogoapi "github.com/Sennue/gogoapi"
)

const (
	READ_BUFFER_SIZE = 1024 * 1024 // 1 meg
)

type AuthCredentials struct {
	DeviceId string `json:"device_id"`
	Password string `json:"password"`
}

type AuthUserInfo struct {
	AccountId int    `json:"account_id"`
	Username  string `json:"username"`
}

type AuthValidator struct {
	db            *sql.DB
	authStatement *sql.Stmt
}

func NewAuthValidator(db *sql.DB) *AuthValidator {
	authStatement, err := db.Prepare(
		"SELECT account_id, name FROM account " +
			"JOIN device USING(account_id) " +
			"WHERE device_id=$1 AND " +
			"password=crypt($2, salt) AND " +
			"active=true;",
	)
	fatal(err)
	return &AuthValidator{db, authStatement}
}

func InternalServerError(err error) gogoapi.JSONError {
	status := http.StatusInternalServerError
	message := fmt.Sprintf("Internal server error.  %s", err.Error())
	return gogoapi.JSONError{status, message}
}

func (validator *AuthValidator) Validate(request *http.Request) (bool, map[string]interface{}, gogoapi.StatusResponse) {
	var credentials AuthCredentials

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	if nil != err {
		return false, nil, InternalServerError(err)
	}
	err = request.Body.Close()
	if nil != err {
		return false, nil, InternalServerError(err)
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		return false, nil, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}
	}

	rows, err := validator.authStatement.Query(credentials.DeviceId, credentials.Password)
	if nil != err {
		return false, nil, InternalServerError(err)
	}
	if !rows.Next() {
		return false, nil, gogoapi.JSONError{http.StatusForbidden, "Authentication failed."}
	} else {
		var (
			account_id int
			username   string
		)
		err = rows.Scan(&account_id, &username)
		if nil != err {
			return false, nil, InternalServerError(err)
		}
		claims := make(map[string]interface{})
		claims["AccessToken"] = "level 1"
		claims["CustomUserInfo"] = AuthUserInfo{
			AccountId: account_id,
			Username:  username,
		}
		return true, claims, nil
	}
}
