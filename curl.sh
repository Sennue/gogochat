#!/bin/sh

export API_HOST='http://localhost:8080'
export CURL_OPTIONS='-sS'
. secure_curl.sh

export API_VERB="GET"
export API_ENDPOINT=""
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="version"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="ping"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="time"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="account"
export API_HEADER="Content-Type: application/json"
export API_DATA="bad-data"
printf "\n$API_VERB /$API_ENDPOINT [bad data]\n"
curl $CURL_OPTIONS -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="account"
export API_HEADER="Content-Type: application/json"
export API_DATA="{\"device_id\":\"$DEVICE_ID\", \"name\":\"$NAME\", \"email\":\"$EMAIL\", \"password\":\"$PASSWORD\"}"
printf "\n$API_VERB /$API_ENDPOINT [may exist]\n"
curl $CURL_OPTIONS -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="account"
export API_HEADER="Content-Type: application/json"
export API_DATA_RANDOM_ID=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f0-9' | head -c 32`
export API_DATA="{\"device_id\":\"FakeID:$API_DATA_RANDOM_ID\", \"name\":\"$API_DATA_RANDOM_ID\", \"email\":\"${API_DATA_RANDOM_ID}@email.com\", \"password\":\"$API_DATA_RANDOM_ID\"}"
printf "\n$API_VERB /$API_ENDPOINT [random values]\n"
curl $CURL_OPTIONS -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="auth"
export API_HEADER="Content-Type: application/json"
export API_DATA="{\"device_id\":\"$DEVICE_ID\", \"password\":\"bad-password\"}"
printf "\n$API_VERB /$API_ENDPOINT [wrong password]\n"
curl $CURL_OPTIONS -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="auth"
export API_HEADER="Content-Type: application/json"
export API_DATA="{\"device_id\":\"$DEVICE_ID\", \"password\":\"$PASSWORD\"}"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="room"
export API_AUTH_HEADER="Authorization: bad-token"
printf "\n$API_VERB /$API_ENDPOINT [bad token]\n"
curl $CURL_OPTIONS -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="room"
export API_AUTH_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="x"
export API_ENDPOINT="room/$API_ITEM"
export API_AUTH_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT [bad syntax]\n"
curl $CURL_OPTIONS -D - -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" #| python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="-1"
export API_ENDPOINT="room/$API_ITEM"
export API_AUTH_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT [nonexistent item]\n"
curl $CURL_OPTIONS -D - -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" #| python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="1"
export API_ENDPOINT="room/$API_ITEM"
export API_AUTH_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="room"
export API_AUTH="Authorization: $TOKEN"
export API_AUTH_HEADER="Content-Type: application/json"
export API_DATA_NAME="Test Room"
export API_DATA_DESCRIPTION="Test room description."
export API_DATA="{ \"name\": \"$API_DATA_NAME\", \"description\": \"$API_DATA_DESCRIPTION\" }"
printf "\n$API_VERB /$API_ENDPOINT [may exist]\n"
curl $CURL_OPTIONS -H "$API_AUTH" -H "$API_AUTH_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="room"
export API_AUTH_HEADER="Authorization: $TOKEN"
export API_HEADER="Content-Type: application/json"
export API_DATA_NAME=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f0-9' | head -c 32`
export API_DATA_DESCRIPTION=`cat /dev/urandom | env LC_CTYPE=C tr -cd 'a-f0-9' | head -c 96`
export API_DATA="{ \"name\": \"$API_DATA_NAME\", \"description\": \"$API_DATA_DESCRIPTION\" }"
printf "\n$API_VERB /$API_ENDPOINT [random values]\n"
curl $CURL_OPTIONS -H "$API_AUTH_HEADER" -H "$API_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="1"
export API_ENDPOINT="message/$API_ITEM"
export API_AUTH_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_AUTH_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="POST"
export API_ENDPOINT="message"
export API_AUTH="Authorization: $TOKEN"
export API_AUTH_HEADER="Content-Type: application/json"
export API_DATA_ROOM_ID="1"
export API_DATA_BODY="Test Message: Hello, World!"
export API_DATA="{ \"room_id\": \"$API_DATA_ROOM_ID\", \"body\": \"$API_DATA_BODY\" }"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_AUTH" -H "$API_AUTH_HEADER" -d "$API_DATA" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" #| python -m json.tool #| less

