.PHONY: all
all: test code_qa

.PHONY: test
test:
	go test
	cd deprecated/cmdline && go test
	cd deprecated/loader && go test
	cd deprecated/validate && go test
	cd deprecated/stringsx && go test
	cd deprecated/sysx && go test
	
.PHONY: code_qa
code_qa:
	go test `go list ./... | grep -v 'hack\|deprecated'` -coverprofile=coverage.txt -covermode=atomic
	golangci-lint run > lint.txt