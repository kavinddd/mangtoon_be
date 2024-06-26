package db

import (
	"database/sql"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)

	}
	testQueries = New(conn)

	os.Exit(m.Run())
}
