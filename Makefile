PORT=8080
DB_DOCKER_CONTAINER=test_container
DB_NAME=unit_testing_db
COVER_OUT=coverage.out
# DSN="host=localhost port=5432 user=root password=secret dbname=${DB_NAME} sslmode=disable timezone=UTC connect_timeout=5"

create_coverage:
	go test ./cmd/web  -coverprofile=${COVER_OUT}

remove_coveragefile:
	rm ${COVER_OUT}

show_coverage:
	go tool cover -html=${COVER_OUT}

test:
	go test -v ./...

postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
# creates the db withing the postgres container
createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=root --owner=root ${DB_NAME}

stop_containers:
	@echo "Stoping all docker containers..."
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers..."; \
		docker stop $$(docker ps -q); \
	else \
		echo "no active containers found..."; \
	fi

start-docker:
	@echo "Running docker container"
	docker start ${DB_DOCKER_CONTAINER}

run: stop_containers start-docker
	go run ./cmd/web
