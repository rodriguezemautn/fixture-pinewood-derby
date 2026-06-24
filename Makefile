.PHONY: init dev build test vet clean all

# ─── Backend ───────────────────────────────────────────────
BACKEND_DIR = backend
BIN_DIR = $(BACKEND_DIR)/bin

init: init-backend init-frontend

init-backend:
	@echo "📦 Initializing backend..."
	@cd $(BACKEND_DIR) && go mod tidy

init-frontend:
	@echo "📦 Initializing frontend..."
	@cd frontend && npm install

dev: kill-backend dev-backend dev-frontend

kill-backend:
	@echo "🔪 Liberando puerto 8080..."
	@lsof -ti:8080 2>/dev/null | xargs -r kill -9 2>/dev/null || true

dev-backend:
	@echo "🚀 Starting backend..."
	@cd $(BACKEND_DIR) && air \
		--build.cmd "go build -o ./tmp/main ./cmd/api" \
		--build.bin "./tmp/main" \
		--build.delay "1" \
		--build.include_dir "cmd internal" \
		--build.exclude_dir "tmp" \
		--build.stop_on_error "true" \
		--log.main_only "true" 2>&1 | sed 's/^/[backend] /'

dev-frontend:
	@echo "🎨 Starting frontend..."
	@cd frontend && npm run dev -- --host 2>&1 | sed 's/^/[frontend] /'

build: build-backend build-frontend

build-backend:
	@echo "🔨 Building backend..."
	@mkdir -p $(BIN_DIR)
	@cd $(BACKEND_DIR) && CGO_ENABLED=0 go build -o bin/api ./cmd/api
	@echo "✅ Backend built: $(BIN_DIR)/api"

build-frontend:
	@echo "🔨 Building frontend..."
	@cd frontend && npm run build
	@echo "✅ Frontend built: frontend/build/"

test:
	@echo "🧪 Running tests..."
	@cd $(BACKEND_DIR) && go test -v -cover ./...

vet:
	@echo "🔍 Running go vet..."
	@cd $(BACKEND_DIR) && go vet ./...

fmt:
	@echo "✨ Checking formatting..."
	@cd $(BACKEND_DIR) && test -z "$$(gofmt -l .)" || (echo "Unformatted files:"; gofmt -l .; exit 1)

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -rf backend/tmp
	@cd frontend && rm -rf build .svelte-kit

# ─── All ───────────────────────────────────────────────────
all: clean init vet test build
	@echo "🎉 All checks passed!"
