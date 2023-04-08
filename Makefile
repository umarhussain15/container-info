
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

CUR_DIR?=$(shell pwd)
OUT_D?=$(shell pwd)/builds
.PHONY: lint
lint:
	@echo "--- Running linter on code ---"
	@golangci-lint run ./... -v
	@echo "--- Running 'go vet' ---"
	@go vet ./...
	@echo "DONE"

.PHONY: test_unit
test_unit:
	@echo "--- Run unit tests ---"
	@go test -v -cover ./...
	@echo "DONE"


.PHONY: build-snapshot
build-snapshot:
	@echo "--- Building app via goreleaser ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) goreleaser build --snapshot --clean
	@echo "built ${OUT_D}/app"
	@echo "DONE"

helm-install:
	@helm install container-info build/helm/container-info