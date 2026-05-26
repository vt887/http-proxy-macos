SHELL := /bin/bash

.PHONY: daemon-dev daemon-test daemon-lint app-open helper-build test format format-check ci security security-govuln security-secrets run-dev db-init render-configs validate-configs

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
	@cd helper && gofmt -w $$(find . -name '*.go')
	@cd app/MacNetGateway && swift format . >/dev/null 2>&1 || true

format-check:
	@UNFORMATTED=$$(find daemon helper -name '*.go' -print0 | xargs -0 gofmt -l); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "Unformatted Go files:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi

ci: format-check test

security: security-govuln security-secrets

security-govuln:
	@GOBIN_PATH=$$(go env GOPATH)/bin; \
	if [ ! -x "$$GOBIN_PATH/govulncheck" ]; then \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi; \
	cd daemon && "$$GOBIN_PATH/govulncheck" ./...; \
	cd ../helper && "$$GOBIN_PATH/govulncheck" ./...

security-secrets:
	@if command -v gitleaks >/dev/null 2>&1; then \
		gitleaks detect --source .; \
	else \
		echo "gitleaks is not installed (install it to run local secret scan)"; \
	fi

run-dev: daemon-dev

db-init:
	@cd daemon && go run ./cmd/macnet-gatewayd >/dev/null 2>&1 & sleep 1; kill $$! || true

render-configs:
	@echo "render-configs is planned for PR-3/PR-4"

validate-configs:
	@echo "validate-configs is planned for PR-3/PR-4"
