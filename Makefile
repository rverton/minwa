generate/sqlc:
	cd internal/database && sqlc generate

generate/templ:
	templ generate

watch:
	templ generate --watch --cmd 'go run ./cmd/minwa'

run:
	go run ./cmd/minwa/
