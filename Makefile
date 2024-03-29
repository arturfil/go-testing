PORT=8080
PASSWORD=secret
USER=root
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

server-coverage:
	go test ./cmd/server -coverprofile=coverage.out && go tool cover -html=coverage.out

test-all:
	go test -v ./...

test-repo:
## go test -v ./pkg/repository/dbrepo
	go test -v -tags=integration ./...

test-server:
	go test -v ./cmd/server 

postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d postgres:12-alpine
# creates the db withing the postgres container
createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${USER} --owner=${USER} ${DB_NAME}

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

run-server: stop_containers start-docker
	go run ./cmd/server/

run-web: stop_containers start-docker
	go run ./cmd/web

expired:
	go run ./cmd/cli/ -action=expired

valid:
	go run ./cmd/cli/ -action=valid
