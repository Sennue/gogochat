# gogochat
JSON chat web API written in Go to demonstrate using the gogoapi micro framework.

## Installation ##
```sh
  # run these commands once
  go get github.com/Sennue/gogochat
  cd $GOPATH/src/github.com/Sennue/gogochat
  ./keys/generate.sh
  cp secure_curl.sh.example secure_curl.sh
  vim secure_curl.sh
  # enter username and password
  go run *.go
  # in a different shell
  ./curl.sh
  # copy token
  vim secure_curl.sh
  # enter token, repeat above steps if it expires
```

## Working With GitHub Keys ##

```sh
  # create ssh key, if need be
  # Reference: https://help.github.com/articles/generating-ssh-keys/
  ssh-keygen -t rsa -C "your_email@example.com"

  # add ssh key to session using the sh shell
  sh
  eval "$(ssh-agent -s)"
  ssh-add ~/.ssh/id_rsa
```
