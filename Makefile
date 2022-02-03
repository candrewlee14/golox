run:
	go run main.go

test:
	go test -tags=unit -cover ./...

cover:
	go test -tags=unit -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v "100" ;
