generate/sqlc:
	cd database && sqlc generate

generate/templ:
	templ generate

run:
	go run ./cmd/minwa/
