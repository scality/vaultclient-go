# Vault Client Golang Module

[![Documentation][badgegodoc]](https://pkg.go.dev/github.com/scality/vaultclient-go/vaultclient)

## Unit Tests
```bash
go test -mod=vendor -v -cover --count 1 ./...
```

## Local Tests

1. Run Vault2
`VAULT_DB_BACKEND=MONGODB yarn dev`
2. Run tests
`go run local-tests/main.go`

## Documentation

[package vaultclient](https://pkg.go.dev/github.com/scality/vaultclient-go/vaultclient)

[badgegodoc]: https://godoc.org/github.com/scality/vaultclient-go/vaultclient?status.svg
