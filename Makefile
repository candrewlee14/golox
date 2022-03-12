run:
	go run main.go

test:
	echo "--- Test Results ---"; \
	gotest -tags=unit,integration -cover ./...;\

cover:
	echo "--- Test Results ---"; \
	gotest -tags=unit,integration -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v \
	-e "100" -e "expressionNode" -e "statementNode"; \
