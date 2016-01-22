PREFIX?=/build

GOFILES = $(shell find . -type f -name '*.go')
apachebeat: $(GOFILES)
	go build

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm elasticbeat || true
