run:
	go run main.go

test:
	echo "--- Test Results ---"; \
	gotest -tags=unit -cover ./...;\

cover:
	echo "--- Test Results ---"; \
	gotest -tags=unit -coverprofile=coverage.out ./... ; \
	go tool cover -func=coverage.out | grep -v \
	-e "100" -e "expressionNode" -e "statementNode"; \
