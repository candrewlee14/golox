run:
	go run main.go

test:
	echo "--- Test Results ---"; \
	go test -tags=unit,integration -coverprofile=coverage.out ./... ; \

cover:
	echo "--- Test Results ---"; \
	go test -tags=unit,integration -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v \
	-e "100" -e "expressionNode" -e "statementNode"; \
	go tool cover -html=coverage.out; \
