db_login:
	psql ${DATABASE_URL}

db_creation_migrate:
	migrate create -ext sql -dir migrations -seq $(name) # make db_creation_migrate name=init_schema

db_migrate:
	migrate -database ${DATABASE_URL} -path migrations up

