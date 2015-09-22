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
printf "\n$API_VERB /$API_ENDPOINT\n"
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
export API_ENDPOINT="rooms"
export API_HEADER="Authorization: bad-token"
printf "\n$API_VERB /$API_ENDPOINT [bad token]\n"
curl $CURL_OPTIONS -H "$API_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ENDPOINT="rooms"
export API_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="x"
export API_ENDPOINT="room/$API_ITEM"
export API_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT [bad syntax]\n"
curl $CURL_OPTIONS -H "$API_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="-1"
export API_ENDPOINT="room/$API_ITEM"
export API_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT [nonexistent item]\n"
curl $CURL_OPTIONS -H "$API_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

export API_VERB="GET"
export API_ITEM="1"
export API_ENDPOINT="room/$API_ITEM"
export API_HEADER="Authorization: $TOKEN"
printf "\n$API_VERB /$API_ENDPOINT\n"
curl $CURL_OPTIONS -H "$API_HEADER" "$API_HOST/$API_ENDPOINT" -X "$API_VERB" | python -m json.tool #| less

