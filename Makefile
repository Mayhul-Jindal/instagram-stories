migrate-init:
	migrate create -ext sql -dir postgres/migration -seq init_schema

migrate-up:
	migrate -path postgres/migration -database "postgresql://postgres:postgres@localhost:5432/instagram-stories?sslmode=disable" -verbose  up 

migrate-down:
	migrate -path postgres/migration -database "postgresql://postgres:postgres@localhost:5432/instagram-stories?sslmode=disable" -verbose  down 