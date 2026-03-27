.PHONY: build build-backend build-frontend test test-backend test-frontend typecheck lint

build: build-backend build-frontend

build-backend:
	@$(MAKE) -C backend build

build-frontend:
	@npm --prefix frontend run build

test: test-backend test-frontend

test-backend:
	@$(MAKE) -C backend test

test-frontend:
	@npm --prefix frontend run typecheck

lint:
	@npm --prefix frontend run lint

typecheck:
	@npm --prefix frontend run typecheck
