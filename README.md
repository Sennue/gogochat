# gogochat
JSON chat web API written in Go to demonstrate using the gogoapi micro framework.

## Global FreeBSD Service Installation ##
```sh
  go get github.com/Sennue/gogochat
  cd $GOPATH/src/github.com/Sennue/gogochat
  go install
  su
  ln -s $GOPATH/bin/gogochat /usr/local/bin/gogochat
  cp gogochat.rc.example /usr/local/etc/rc.d/gogochat
  mkdir /usr/local/etc/gogochat
  mkdir /usr/local/etc/gogochat/keys
  cd /usr/local/etc/gogochat/keys/
  openssl genrsa -out app.rsa 4096
  openssl rsa -in app.rsa -pubout > app.rsa.pub
  pw adduser gogochat -g gogochat -d /nonexistent -s /usr/sbin/nologin -c "gogochat service user"
  vim /etc/rc.conf
  # add gogochat_enable="YES" to run at boot
  # Use the following to control the service:
  # service gogochat start/restart/stop/status
  # Use the following to control the service without starting at boo:
  # service gogochat onestart/onerestart/onestop/onestatus
```

## Local Development Installation ##
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
