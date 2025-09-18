

go install github.com/go-task/task/v3/cmd/task@latest

GOBIN=./bin go install github.com/go-task/task/v3/cmd/task@latest



// protobuf
brew install protobuf


// grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


// swagger
go install \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
api/service.proto