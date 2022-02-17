# This file will use so our app will wait from the db to be up before trying to connect

# docker enirement vars
wait-for "${DATABASE_HOST}:${DATABASE_PORT}" -- "$@"

# watch for .go file changes
CompileDaemon --build="go build ./cmd/web" --command=./web