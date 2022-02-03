run:
	go run main.go

test:
	go test -cover ./...

cover:
	go test -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v "100"; \
