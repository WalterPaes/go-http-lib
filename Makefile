tests:
	go run test ./...
cover:
	go test -coverprofile=coverage.out ./...
cover-html:
	go tool cover -html=coverage.out
godoc:
	go doc --all ./http/