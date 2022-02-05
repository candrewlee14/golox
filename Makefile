run:
	go run main.go

test:
	gotest -tags=unit -cover ./...

cover:
	gotest -tags=unit -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v "100" ;
