.PHONY: all
all: test

.PHONY: test
test:
	go test
	cd stdlibx/cmdline && go test
	cd stdlibx/loader && go test
	cd stdlibx/settings && go test
	cd stdlibx/validate && go test
	cd stdlibx/stringsx && go test
	cd stdlibx/sysx && go test
	# remove in the future
	cd deprecated/provider && go test
	
.PHONY: test_coverage
test_coverage:
	go test `go list ./... | grep -v 'hack\|deprecated'` -coverprofile=coverage.txt -covermode=atomic
