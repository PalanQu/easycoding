# https://betterprogramming.pub/my-ultimate-makefile-for-golang-projects-fcc8ca20c9bb

# TODO(qujiabao): add args to these command
# docker-build
# docker-release
# setup

# deps
.PHONY: deps
deps:
	go mod tidy


# manage
.PHONY: gen-config
gen-config:
	go run cmd/manage/main.go config generate

.PHONY: gen-api
gen-api:
	go run cmd/manage/main.go proto gen-go
	go run cmd/manage/main.go proto gen-swagger

.PHONY: clean
clean:
	go run cmd/manage/main.go proto clean-go
	go run cmd/manage/main.go proto clean-swagger


# build and run
.PHONY: build
build: gen-api
	go build cmd/serve/main.go

.PHONY: run
run: gen-api
	docker-compose up mysql -d
	sleep 2
	EASYCODING_DATABASE_HOST=localhost go run cmd/serve/main.go

.PHONY: stop
stop:
	docker-compose down

.PHONY: run-in-docker
run-in-docker:
	docker-compose up


# test
.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -cover ./...

.PHONY: coverage-html
coverage-html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out


# create migrate
.PHONY: migrate-create
migrate-create:
	go run cmd/migrate/main.go create --all


# lint
.PHONY: lint
lint:
	pre-commit run --all
