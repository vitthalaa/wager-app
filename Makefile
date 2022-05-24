test:
	@go test -race $(shell go list ./... | grep -v /vendor/) -v -coverprofile=coverage.txt

integration-test:
	@go test ./integration_tests/ -tags=integration

gen:
	@go generate ./...
	@goimports -local github.com/vitthalaa/wager-app -w .

run:
	@go run main.go

build:
	@GOOS=linux GOARCH=amd64 go build -o wager-app main.go

docker-verify:
	@docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit --force-recreate

docker-run:
	@docker-compose up -d

docker-down:
	@docker-compose down



