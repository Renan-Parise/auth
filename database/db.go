package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/Renan-Parise/codium/errors"
	"github.com/Renan-Parise/codium/utils"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDBInstance() *sql.DB {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		database, err := sql.Open("mysql", dsn)
		if err != nil {
			utils.GetLogger().WithError(err).Error("Error connecting to database")

			errors.NewDatabaseError(err.Error())
			panic(err)
		}
		db = database
	})
	return db
}
