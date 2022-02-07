.PHONY: all
all: test

.PHONY: test
test:
	cd env && go test
	cd id && go test
	cd loader && go test
	cd settings && go test
	cd timestamp && go test
	cd validate && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|google'` -coverprofile=coverage.txt -covermode=atomic
