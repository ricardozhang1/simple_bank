package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/techschool/simplebank/util"
	"log"
	"os"
	"testing"
)

// Queries 对象是用于操作db的句柄， 可对*sql.DB对象使用New方法得到
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configure: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}