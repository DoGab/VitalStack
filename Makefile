# MacroGuard Development

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
	@echo "MacroGuard Development Commands:"
	@echo ""
	@echo "  make dev         - Start both backend and frontend (recommended)"
	@echo "  make dev-api     - Start Go API server only (with Air hot reload)"
	@echo "  make dev-web     - Start SvelteKit frontend only"
	@echo "  make install     - Install all dependencies"
	@echo "  make install-api - Install Go dependencies"
	@echo "  make install-web - Install frontend dependencies"
	@echo "  make clean       - Remove build artifacts"
	@echo ""
	@echo "Endpoints:"
	@echo "  Frontend:  http://localhost:5173"
	@echo "  API:       http://localhost:8080"
	@echo "  API Docs:  http://localhost:8080/docs"
