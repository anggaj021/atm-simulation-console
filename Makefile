run:
	go run app/main.go
test:
	go test -v ./... -coverprofile=cov.out
coverage:
	go tool cover -html=cov.out