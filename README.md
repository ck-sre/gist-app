# gist-app
Gist app application in Go. Presents gists that are updated from by users. Also has support for login. Backend is a local Mysql.

Run like this:
go run ./cmd/http-web -port=":9999" -dbconn='<usr>:<pwd>@/gistapp?parseTime=true'

