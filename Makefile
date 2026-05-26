SHELL := /bin/bash

.PHONY: daemon-dev daemon-test daemon-lint app-open helper-build test format run-dev db-init render-configs validate-configs

daemon-dev:
	@cd daemon && go run ./cmd/macnet-gatewayd

daemon-test:
	@cd daemon && go test ./...

daemon-lint:
	@echo "lint tool is not configured in PR-1"

app-open:
	@open app/MacNetGateway

helper-build:
	@cd helper && go build ./cmd/macnet-helper

test: daemon-test
	@cd app/MacNetGateway && swift test

format:
	@cd daemon && gofmt -w $$(find . -name '*.go')
	@cd app/MacNetGateway && swift format . >/dev/null 2>&1 || true

run-dev: daemon-dev

db-init:
	@cd daemon && go run ./cmd/macnet-gatewayd >/dev/null 2>&1 & sleep 1; kill $$! || true

render-configs:
	@echo "render-configs is planned for PR-3/PR-4"

validate-configs:
	@echo "validate-configs is planned for PR-3/PR-4"
