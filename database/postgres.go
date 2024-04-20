package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	zerolog "github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

var (
	DB *sqlx.DB
)

// ConnectDB function: Make database connection
func init() {

	// Load environmenatal variables
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DB = sqlx.NewDb(SqlDb(), "postgres")
}

func SqlDb() *sql.DB {
	dsn := DbConnectionString()
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalln(err)
	}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zlogger := zerolog.New(os.Stdout).With().Logger()
	// prepare logger
	loggerOptions := []sqldblogger.Option{
		sqldblogger.WithSQLQueryFieldname("sql"),
		// sqldblogger.WithWrapResult(false),
		/* 		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		   		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		   		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug), */
	}
	db = sqldblogger.OpenDriver(dsn, db.Driver(), zerologadapter.New(zlogger), loggerOptions...)
	return db
}

func DbConnectionString() string {

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DB")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, databaseHost, databasePort, databaseName)
	return url

}

func HasUniquenessViolation(err error) (bool, string) {

	if err, ok := err.(*pq.Error); ok {
		return err.Code.Name() == "unique_violation", err.Constraint
	}

	return false, ""
}

func ErrorKey(err error) string {
	switch err {
	case sql.ErrNoRows:
		return "not_found"
	default:
		return "unknown_db_error"
	}
}
