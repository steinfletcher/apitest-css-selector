.PHONY: test test-examples docs fmt vet

test:
	bash -c 'diff -u <(echo -n) <(gofmt -s -d ./)'
	bash -c 'diff -u <(echo -n) <(go vet ./...)'
	go test ./...
