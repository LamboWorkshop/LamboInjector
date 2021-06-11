package cmd

// import (
// 	"database/sql"
// 	"fmt"
// )

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "awesome-pass"
// 	// dbname   = "postgres"
// 	dbname = "lambo_bot"
// )

// type Database struct {
// 	Host     string
// 	Port     int
// 	User     string
// 	Password string
// 	Dbname   string
// }

// func connectToDB() (*sql.DB, error) {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)
// 	// os.Getenv("DB_HOST"), port, user, password, dbname)

// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		return db, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		return db, err
// 	}

// 	return db, err
// }

// func createTable(db *sql.DB, tabeName string) error {

// 	createTable := `
// 	CREATE TABLE ` + tabeName + ` (
// 	id  SERIAL PRIMARY KEY,
// 	timestamp BIGINT,
// 	amount TEXT,
// );
// `
// 	_, err := db.Query("select * from fiat_deposit;")
// 	if err == nil {
// 		return nil
// 	}

// 	stmt, err := db.Prepare(createTable)
// 	defer stmt.Close()

// 	_, err = stmt.Exec()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func InsertToDB(db *sql.DB, table string) error {
// 	id := 0
// 	sqlStatement := `
// 	INSERT INTO ` + table + ` (timestamp, amount)
// 	VALUES ($1, $2)
// 	RETURNING id`

// 	return db.QueryRow(sqlStatement, timestamp, amount).Scan(&id)
// }
