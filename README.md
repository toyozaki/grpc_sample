# grpc_sample

# Prerequisite
```
$ brew install protobuf
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# How to generate grpc-code
```
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. \
       --go-grpc_out=. \
       --go_opt=paths=source_relative \
       --go-grpc_opt=paths=source_relative \
       gen/greet.proto
```

# How to run server
```
go run server/main.go
```