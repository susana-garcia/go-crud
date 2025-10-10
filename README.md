
Go project with a server that serves gRPC(s) and uses a PostgreSQL database.

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
