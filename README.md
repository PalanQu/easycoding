# Easy coding

This repo contains an example structure for a monolithic Go Web Application.

## Project Architecture

This project loosely follows [Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

<img src="docs/pics/project_architecture.jpg" style="zoom:50%;" />

## Features

- 100% API defined by protobuf
- Auto generate grpc, grpc gateway, validate go files
- Provide both rest api and grpc api
- Auto generate swagger api document
- Builtin prometheus metrics
- Support import api definition by postman
- Run in docker
- Auto configuration generate
- Database migrate up and down
- Database mock testing
- Golang, Protobuf and basic text file linting
- Error definition and classification
- Auto logging and pretty format
- Unit test and test coverage
- Graceful stop
- Backend processes
- Health check

## TODO

- Use reflect in configration
- Benchmark
- Fix linting
- Intergration test
- Auth
- More options in configuration
- Property based test
- GraphQL server
