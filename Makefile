.PHONY: all
all: test

.PHONY: test
test:
	cd pkg/env && go test
	cd pkg/id && go test
	cd pkg/timestamp && go test
	cd pkg/validate && go test
	cd pkg/loader && go test
	cd pkg/observer && go test
	cd pkg/observer/provider && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'tests\|google\|httpserver'` -coverprofile=coverage.txt -covermode=atomic
