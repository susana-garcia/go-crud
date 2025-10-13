
# Go service

Go project with a server that serves gRPC(s) request for microservices communication and uses a PostgreSQL database.

## Getting started

1. Create `.env.local` with the environments variables
1. Run `make run`

## Database

A postgres database is required to run the application. To start a docker image, you can run `make start-postgres` and to stop it, you can run `make stop-postgres`.

You can connect to it using any postgres client: 
  - User: `gocrud`
  - Password: `gocrud`
  - Host: `localhost`
  - Port: `5435`
  - DB: `gocrud`

## Protobuf generation

Buf-based workflow that manages dependencies:

- `buf.yaml` - Defines protobuf module and dependencies.
- `buf.gen.yaml` - Configures code generation.
- `scripts/generate-pb.sh` - Uses buf directly.

Whenever changes are made in `pb/blog.proto`, we will need to rerun `make generate`.

## Run locally

Once a local database is up and running, you can run the scripts to create, get, update or delete.

## Tests

To run tests:
1. `make start-postgres-test`
1. `make test` 

## Helpful resources

- [GORM guide](https://gorm.io/docs/index.html)
- [Protovalidate docs](https://protovalidate.com/quickstart/grpc-go/)
- [gRPC docs](https://grpc.io/docs/languages/go/quickstart/)
- [grpcurl](https://github.com/fullstorydev/grpcurl)
