package storage

import (
	"database/sql"
	"github.com/google/logger"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var Db *sql.DB
var GormDB *gorm.DB

func Init(connStr string) error {
	var err error
	Db, err = sql.Open("postgres", connStr)
	GormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: Db,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	return err
}

func NewListener(connStr string) *pq.Listener {
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logger.Fatalf("Postgres listener failure: %v", err)
		}
	}
	minReconn := 10 * time.Second
	maxReconn := time.Minute
	return pq.NewListener(connStr, minReconn, maxReconn, reportProblem)
}
