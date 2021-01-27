test:
	@ go test -race -v ./...

lint:
	@ golangci-lint run \
			--disable-all \
			--enable goimports \
			--enable govet \
			--enable misspell \
			--enable errcheck \
			--enable staticcheck \
			--enable prealloc \
			--enable ineffassign

gen:
	@ go generate ./...