#!/bin/sh

KEYSIZE=4096

openssl genrsa -out app.rsa $KEYSIZE
openssl rsa -in app.rsa -pubout > app.rsa.pub

