module example/goWorld

go 1.26.2

require (
	connectrpc.com/connect v1.21.0
	github.com/jackc/pgx/v5 v5.11.0
	google.golang.org/protobuf v1.36.12
	rsc.io/quote v1.5.2
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	rsc.io/sampler v1.3.0 // indirect
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	google.golang.org/protobuf/cmd/protoc-gen-go
)
