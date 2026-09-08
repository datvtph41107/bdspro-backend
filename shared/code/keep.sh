protoc --go_out=. --go_opt=paths=source_relative chat.proto
swag init --parseDependency
kill -9 $(lsof -ti tcp:8032)
