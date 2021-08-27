.PHONY: all
all: test

.PHONY: test
test:
	cd pkg/env && go test
	cd pkg/id && go test
	cd pkg/loader && go test
	cd pkg/provider && go test
	cd pkg/timestamp && go test
	cd pkg/validate && go test
	cd observer && go test
	cd storage && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|google'` -coverprofile=coverage.txt -covermode=atomic
