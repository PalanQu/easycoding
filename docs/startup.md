# startup

1. install protoc

protoc https://github.com/protocolbuffers/protobuf/releases

2. install protoc plugin

``` bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest \
    google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    github.com/envoyproxy/protoc-gen-validate@latest
```

3. install protobuf management

buf https://github.com/bufbuild/buf/releases

4. install go-swagger cli

https://github.com/go-swagger/go-swagger/releases

7. docker

8. docker-compose

``` bash
sudo curl -L "https://github.com/docker/compose/releases/download/v2.6.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

```

- upgrade python to 3.8

5. install precommit

``` bash
pip3 install pre-commit
pre-commit install
```

6. golang lint

``` bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
