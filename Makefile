run-docker-stdout:
	docker compose up -d postgres-auth
	docker compose build
	docker compose up go-auth

run-docker-background:
	docker compose up -d postgres-auth
	docker compose build
	docker compose up -d go-auth

stop-docker:
	docker compose down