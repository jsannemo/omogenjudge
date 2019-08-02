package db

import (
	"fmt"
	"flag"
	"strings"
  "time"

	"github.com/google/logger"
	"github.com/lib/pq"
  "github.com/jmoiron/sqlx"
)

var (
	dbName     = flag.String("db_name", "omogenjudge", "Name of the database used for application storage")
	dbUser     = flag.String("db_user", "omogenjudge", "Name of the user that should connect to the database")
	dbPassword = flag.String("db_password", "omogenjudge", "Password used to connect to the database")
	dbHost     = flag.String("db_host", "localhost:5432", "Host in the form host:port that the database listens to")
)

var pool *sqlx.DB // Database connection pool.

func connString() string {
	hostPort := strings.Split(*dbHost, ":")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    hostPort[0], hostPort[1], *dbUser, *dbPassword, *dbName)
}

func Init() {
	pool = sqlx.MustConnect("postgres", connString())
	logger.Infof("Connected to database: %v", *dbName)
}

func Conn() (*sqlx.DB) {
  if pool == nil {
    Init()
  }
  return pool
}

func NewListener() *pq.Listener {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
      logger.Fatalf("Postgres listener failure: %v", err)
		}
	}

	minReconn := 10 * time.Second
	maxReconn := time.Minute
  return pq.NewListener(connString(), minReconn, maxReconn, reportProblem)
}
