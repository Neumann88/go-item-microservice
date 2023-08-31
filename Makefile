VERSION?=dev

.PHONY:dev
start:
	./scripts/env.sh ${VERSION} "go run ./cmd/*.go"

.PHONY: compose-up
compose-up:
	docker compose up --build && docker-compose logs -f

.PHONY: compose-down
compose-down:
	docker compose down --remove-orphans --volumes

.PHONY: remove-volumes
remove-volumes:
	docker system prune -a --volumes