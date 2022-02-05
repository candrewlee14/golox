run:
	go run main.go

test:
	echo "";
	gotest -tags=unit -cover ./...;\

cover:
	echo "";
	gotest -tags=unit -coverprofile=coverage.out ./... ; \
	echo \n; \
	go tool cover -func=coverage.out | grep -v \
	-e "100" -e "expressionNode" -e "statementNode";
