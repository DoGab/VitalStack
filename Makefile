# MacroGuard Development

# Simple color scheme
CYAN := \033[36m
DIM := \033[2m
RESET := \033[0m

# Basic symbols
CHECK := ✓
ARROW := →

.PHONY: dev dev-api dev-web install install-api install-web clean help

# Start both backend and frontend concurrently
dev:
	@echo "🚀 Starting MacroGuard development servers..."
	@make -j2 dev-api dev-web

# Start Go backend with Air hot reload
dev-api:
	@echo "🔧 Starting Go API server with Air hot reload..."
	@cd apps/api-go && air

# Start SvelteKit frontend dev server
dev-web:
	@echo "🌐 Starting SvelteKit frontend..."
	@cd apps/web && bun run dev

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

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf apps/api-go/tmp
	@rm -rf apps/web/.svelte-kit
	@rm -rf apps/web/node_modules

# Show available commands
help:
	@printf "\n${CYAN}MacroGuard Development Commands${RESET}\n"
	@printf "$(DIM)────────────────────────────────────$(RESET)\n"
	@printf "  make dev         $(ARROW) Start both backend and frontend (recommended)\n"
	@printf "  make dev-api     $(ARROW) Start Go API server only (with Air hot reload)\n"
	@printf "  make dev-web     $(ARROW) Start SvelteKit frontend only\n"
	@printf "  make install     $(ARROW) Install all dependencies\n"
	@printf "  make install-api $(ARROW) Install Go dependencies\n"
	@printf "  make install-web $(ARROW) Install frontend dependencies\n"
	@printf "  make clean       $(ARROW) Remove build artifacts\n"
	@printf "\n"
	@printf "Endpoints:\n"
	@printf "  Frontend:  ${CYAN}http://localhost:5173${RESET}\n"
	@printf "  API:       ${CYAN}http://localhost:8080${RESET}\n"
	@printf "  API Docs:  ${CYAN}http://localhost:8080/docs${RESET}\n"
