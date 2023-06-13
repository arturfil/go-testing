COVER_OUT=coverage.out

create_coverage:
	go test ./cmd/web  -coverprofile=${COVER_OUT}

remove_coveragefile:
	rm ${COVER_OUT}

show_coverage:
	go tool cover -html=${COVER_OUT}

make test:
	go test ./cmd/web -v

run:
	go run ./cmd/web
