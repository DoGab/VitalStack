# MacroGuard Development

# Simple color scheme
CYAN := \033[36m
DIM := \033[2m
RESET := \033[0m

# Basic symbols
CHECK := ✓
ARROW := →

.PHONY: dev dev-api dev-web install install-api install-web clean help \
        lint lint-api lint-web format format-api format-web fix fix-api fix-web \
        openapi gen-openapi-spec gen-api-client

# ─────────────────────────────────────
# Development
# ─────────────────────────────────────

# Generate OpenAPI spec and TypeScript client
openapi: gen-openapi-spec gen-api-client
	@echo "✅ OpenAPI spec and client generated"

# Generate OpenAPI spec from Go API
gen-openapi-spec:
	@echo "📝 Generating OpenAPI spec..."
	@cd apps/api-go && go run main.go openapi --config local-config.yaml

# Generate TypeScript API client from OpenAPI spec
gen-api-client:
	@echo "🔄 Generating TypeScript API client..."
	@cd apps/web && bun run generate:api

# Start both backend and frontend concurrently
dev:
	@echo "🚀 Starting MacroGuard development servers..."
	@make -j2 dev-api dev-web

# Start Go backend with Air hot reload
dev-api:
	@echo "🔧 Starting Go API server with Air hot reload..."
	@cd apps/api-go && air

# Start SvelteKit frontend dev server (network accessible)
dev-web:
	@echo "🌐 Starting SvelteKit frontend..."
	@cd apps/web && bun run dev --host

# ─────────────────────────────────────
# Install
# ─────────────────────────────────────

# Install all dependencies
install: install-api install-web
	@echo "✅ All dependencies installed"

# Install Go dependencies
install-api:
	@echo "📦 Installing Go dependencies..."
	@cd apps/api-go && go mod tidy

# Install frontend dependencies
install-web:
	@echo "📦 Installing frontend dependencies..."
	@cd apps/web && bun install

# ─────────────────────────────────────
# Lint
# ─────────────────────────────────────

# Lint both projects
lint: lint-api lint-web
	@echo "✅ Linting complete"

# Lint Go API
lint-api:
	@echo "🔍 Linting Go API..."
	@cd apps/api-go && golangci-lint run

# Lint frontend
lint-web:
	@echo "🔍 Linting frontend..."
	@cd apps/web && bun run lint

# ─────────────────────────────────────
# Format
# ─────────────────────────────────────

# Format both projects
format: format-api format-web
	@echo "✅ Formatting complete"

# Format Go API
format-api:
	@echo "✨ Formatting Go API..."
	@cd apps/api-go && go fmt ./...

# Format frontend
format-web:
	@echo "✨ Formatting frontend..."
	@cd apps/web && bun run format

# ─────────────────────────────────────
# Lint + Format (fix)
# ─────────────────────────────────────

# Fix all lint and format issues
fix: fix-api fix-web
	@echo "✅ All fixes applied"

# Fix Go API issues
fix-api:
	@echo "🔧 Fixing Go API..."
	@cd apps/api-go && golangci-lint run --fix && go fmt ./...

# Fix frontend issues
fix-web:
	@echo "🔧 Fixing frontend..."
	@cd apps/web && bun run lint:fix && bun run format

# ─────────────────────────────────────
# Clean
# ─────────────────────────────────────

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf apps/api-go/tmp
	@rm -rf apps/web/.svelte-kit
	@rm -rf apps/web/node_modules

# ─────────────────────────────────────
# Help
# ─────────────────────────────────────

# Show available commands
help:
	@printf "\n${CYAN}MacroGuard Development Commands${RESET}\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@printf "\n${CYAN}Development${RESET}\n"
	@printf "  make dev           $(ARROW) Start both backend and frontend\n"
	@printf "  make dev-api       $(ARROW) Start Go API server (with Air)\n"
	@printf "  make dev-web       $(ARROW) Start SvelteKit frontend (network accessible)\n"
	@printf "\n${CYAN}Code Quality${RESET}\n"
	@printf "  make lint          $(ARROW) Lint all projects\n"
	@printf "  make lint-api      $(ARROW) Lint Go API only\n"
	@printf "  make lint-web      $(ARROW) Lint frontend only\n"
	@printf "  make format        $(ARROW) Format all projects\n"
	@printf "  make format-api    $(ARROW) Format Go API only\n"
	@printf "  make format-web    $(ARROW) Format frontend only\n"
	@printf "  make fix           $(ARROW) Fix all lint + format issues\n"
	@printf "  make fix-api       $(ARROW) Fix Go API issues\n"
	@printf "  make fix-web       $(ARROW) Fix frontend issues\n"
	@printf "\n${CYAN}Setup${RESET}\n"
	@printf "  make install       $(ARROW) Install all dependencies\n"
	@printf "  make install-api   $(ARROW) Install Go dependencies\n"
	@printf "  make install-web   $(ARROW) Install frontend dependencies\n"
	@printf "  make clean         $(ARROW) Remove build artifacts\n"
	@printf "\n${CYAN}Code Generation${RESET}\n"
	@printf "  make openapi       $(ARROW) Generate spec + TypeScript client\n"
	@printf "  make gen-openapi-spec $(ARROW) Generate OpenAPI spec only\n"
	@printf "  make gen-api-client $(ARROW) Generate TypeScript client only\n"
	@printf "\n"
	@printf "Endpoints:\n"
	@printf "  Frontend:  ${CYAN}http://localhost:5173${RESET}\n"
	@printf "  API:       ${CYAN}http://localhost:8080${RESET}\n"
	@printf "  API Docs:  ${CYAN}http://localhost:8080/docs${RESET}\n"
