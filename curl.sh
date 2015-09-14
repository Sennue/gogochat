#!/bin/sh

export API_HOST='http://localhost:8080'
. secure_curl.sh

#export API_ENDPOINT=""
#curl "$API_HOST/$API_ENDPOINT" -X GET | python -m json.tool | less

#export API_ENDPOINT="todos"
#curl -H "Authorization: BadToken" "$API_HOST/$API_ENDPOINT" -X GET | python -m json.tool | less

export API_ENDPOINT="auth"
curl -H "Content-Type: application/json" -d "{\"username\":\"$USERNAME\", \"password\":\"$PASSWORD\"}" "$API_HOST/$API_ENDPOINT" -X POST | python -m json.tool | less

export API_ENDPOINT="user"
curl -H "Authorization: $TOKEN" "$API_HOST/$API_ENDPOINT" -X GET | python -m json.tool | less

#export API_ITEM="12345"
#export API_ENDPOINT="todos/$API_ITEM"
#curl -H "Authorization: $TOKEN" "$API_HOST/$API_ENDPOINT" -X GET | python -m json.tool | less

