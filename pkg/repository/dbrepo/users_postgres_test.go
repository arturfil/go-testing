package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"webapp/pkg/data"
	"webapp/pkg/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "user_test"
	port     = "5435"
	dsn      = "host=%s port=%s password=%s user=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.DatabaseRepo

func TestMain(m *testing.M) {
	// connect to docker;
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker %s", err)
	}

	pool = p

	// set up our docker options
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// get a resource
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
        _ = pool.Purge(resource)
	}

	// start the image adn wait until it's ready
	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		// _ = pool.Purge(resource)
		log.Fatalf("could not connect to the database: %s", err)
	}

	// populate the database with empty talbes
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

    testRepo = &PostgresDBRepo{DB: testDB}

	// run tests
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource %s", err)
	}

	// clean up
	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("cannont ping database")
	}
}

func TestPostgresDBRepoInsertUser(t *testing.T) {
    testUser := data.User {
        FirstName: "Admin",
        LastName: "User",
        Email: "admin@example.com",
        Password: "secret",
        IsAdmin: 1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    id, err := testRepo.InsertUser(testUser)
    if err != nil {
        t.Errorf("insert user returned an error: %s", err)
    }

    if id != 1 {
        t.Errorf("inset user retunred wrong id; expected 1, but got %d", id)
    }
}
