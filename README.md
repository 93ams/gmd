# Gommand
Gommand is a simple command line tool that allows you to run commands in a structured way.

## Tasks
### Boot
Bootstraps the project
```sh
go install github.com/joerdav/xc/cmd/xc@latest
go install github.com/charmbracelet/gum@latest
```
### ShowEnv
```sh
go run ./cmd env show
```
### SelectEnv
```sh
go run ./cmd env select $(gum choose "" $(go run ./cmd env list))  
```
### ListEnvs
```sh
go run ./cmd env list
```
### CurrentEnv
```sh
go run ./cmd env current
```
### AddEnv
```sh
go run ./cmd env add
```
### RemEnv
```sh
go run ./cmd env rem
```