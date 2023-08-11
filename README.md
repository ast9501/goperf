

## Build for alpine OS
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build server.go
```
See [this](https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host) post
