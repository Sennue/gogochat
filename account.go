package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
"fmt"

	gogoapi "github.com/Sennue/gogoapi"
)

type AccountObject struct {
	DeviceId string `json:"device_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountResource struct {
	auth            *gogoapi.AuthResource
	db              *sql.DB
	createStatement *sql.Stmt
	//deleteStatement *sql.Stmt
}

func NewAccountResource(auth *gogoapi.AuthResource, db *sql.DB) *AccountResource {
	createStatement, err := db.Prepare(
		"SELECT * FROM add_account($1, $2, $3, $4, true);",
	)
	fatal(err)
	//deleteStatement, err := db.Prepare()
	return &AccountResource{auth, db, createStatement}
}

// TODO: function to add device
func (accountResource *AccountResource) Post(request *http.Request) (int, interface{}, http.Header) {
	var accountObject AccountObject

	body, err := ioutil.ReadAll(io.LimitReader(request.Body, READ_BUFFER_SIZE))
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	err = request.Body.Close()
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}

	if err := json.Unmarshal(body, &accountObject); err != nil {
		return gogoapi.HTTP_UNPROCESSABLE, gogoapi.JSONError{gogoapi.HTTP_UNPROCESSABLE, "Unprocessable entity."}, nil
	}
fmt.Println(accountObject)
	rows, err := accountResource.createStatement.Query(
		accountObject.DeviceId,
		accountObject.Name,
		accountObject.Email,
		accountObject.Password,
	)
	if nil != err {
		return http.StatusInternalServerError, InternalServerError(err), nil
	}
	if !rows.Next() {
		return http.StatusInternalServerError, InternalServerError(err), nil
	} else {
		var (
			success    bool
			account_id int
		)
		err = rows.Scan(&success, &account_id)
		if nil != err {
			return http.StatusInternalServerError, InternalServerError(err), nil
		} else if !success {
			return http.StatusConflict, gogoapi.JSONError{http.StatusForbidden, "Already registered."}, nil
		}
		claims := AuthClaims("level 1", accountObject.Name, account_id)
		return accountResource.auth.AuthTokenResponse(claims)
	}
}
