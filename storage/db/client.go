// Database client and utilities for setting up a connection
package db

import (
	"database/sql"
	"fmt"
	"flag"
	"strings"
  "time"

	"github.com/google/logger"
	"github.com/lib/pq"
)

var (
	dbName     = flag.String("db_name", "omogenjudge", "Name of the database used for application storage")
	dbUser     = flag.String("db_user", "omogenjudge", "Name of the user that should connect to the database")
	dbPassword = flag.String("db_password", "omogenjudge", "Password used to connect to the database")
	dbHost     = flag.String("db_host", "localhost:5432", "Host in the form host:port that the database listens to")
)

var pool *sql.DB // Database connection pool.

func connString() string {
	hostPort := strings.Split(*dbHost, ":")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    hostPort[0], hostPort[1], *dbUser, *dbPassword, *dbName)
}

func init() {
	var err error
	pool, err = sql.Open("postgres", connString())
	if err != nil {
		logger.Fatalf("Could not connect to database: %v", err)
	}

	if err = pool.Ping(); err != nil {
		logger.Fatalf("Could not ping database: %v", err)
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)
	logger.Infof("Connected to database: %v", *dbName)
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

func GetPool() (*sql.DB) {
	return pool
}

func NewListener() *pq.Listener {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
      logger.Errorf("Postgres listener failure: %v", err)
		}
	}

	minReconn := 10 * time.Second
	maxReconn := time.Minute
  return pq.NewListener(connString(), minReconn, maxReconn, reportProblem)
}
