package db

import (
    "database/sql"
    "log"
    "os"
    "testing"
	_"github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

const (
    dbDriver = "postgres"
    dbSource = "postgresql://root:secret@localhost:3000/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
    var err error
    testDb, err = sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("Cannot Connect to the Database:", err)
    }

    testQueries = New(testDb)
    exitCode := m.Run()

   
    os.Exit(exitCode)
}
