
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
goose postgres postgres://username:@localhost:5432/golangWebRss up
goose postgres postgres://username:@localhost:5432/golangWebRss down

* run this from root of the folder - sqlc generate 