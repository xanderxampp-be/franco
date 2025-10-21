package dbmicro

import (
	"database/sql"
	"os"
	"time"

	"github.com/xanderxampp-be/franco/log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/mysql"
)

func New(driverName string) *sql.DB {
	db, err := apmsql.Open(driverName, os.Getenv("USER_DB")+":"+os.Getenv("PASS_DB")+"@tcp("+os.Getenv("HOST_DB")+")/"+os.Getenv("SCHEMA_NAME"))

	if err != nil {
		panic("fail to connect to database")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(40)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

// NewMysql is a
func NewMysql(user, pass, host, dbname string, maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime int) *sql.DB {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	cfg := mysql.Config{
		User:                 user,
		Passwd:               pass,
		Addr:                 host,
		DBName:               dbname,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  loc,
	}

	db, err := apmsql.Open(
		"mysql",
		cfg.FormatDSN(),
	)

	if err != nil {
		log.LogDebug("error connection : " + err.Error())
		panic("fail to connect to database")
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(connMaxIdleTime) * time.Minute)

	return db
}
