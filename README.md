## Overview

This is a library meant to power the ZFlights product. It contains code for the backend for a website and code for the daily lambda refresh.

## Testing

To run the tests, you must start up the db by first running the following command to start up the local database.
```shell
docker compose up -d
```
You can then run the tests using the following `go` command.
```shell
go test ./...
```
