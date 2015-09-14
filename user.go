package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"

	"code.google.com/p/go.crypto/scrypt"
)

const (
	SALT_BYTES = 32
	HASH_BYTES = 64
)

type UserResource struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

// pass "" as salt to generate new salt and hash password
// otherwise password is assumed to already have been hashed with provided salt
func NewUserResource(id int, name, password, salt string) *UserResource {
	if "" == salt {
		saltBytes := make([]byte, SALT_BYTES)
		_, err := io.ReadFull(rand.Reader, saltBytes)
		if err != nil {
			log.Fatal(err)
		}
		salt = hex.EncodeToString(saltBytes)
		password = HashPassword(password, salt)
	}
	return &UserResource{id, name, password, salt}
}

func HashPassword(password, salt string) string {
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		log.Fatal(err)
	}
	passwordBytes, err := scrypt.Key([]byte(password), saltBytes, 1<<14, 8, 1, HASH_BYTES)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(passwordBytes)
}

func (user *UserResource) Get(request *http.Request) (int, interface{}, http.Header) {
	return 200, user, nil
}
