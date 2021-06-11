package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	host     = os.Getenv("POSTGRES_HOST")
	port     = 5432
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

type sCoinsToMonitor struct {
	pair    string
	status  string
	percent string
	id      int
}

var dbClient *sql.DB

func connectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, err
}

func initDataBase() (*sql.DB, error) {
	db, err := connectToDB()
	if err != nil {
		logrus.Println("❌ initDatabase():", err)
	} else {
		logrus.Println("✅ initDatabase()")
	}

	return db, err
}

func createTable(db *sql.DB, tableName, tableContent string) error {

	query := `CREATE TABLE IF NOT EXISTS ` + tableName + ` (` + tableContent + `);`

	stmt, err := db.Prepare(query)
	if err != nil {
		logrus.Error(err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func checkBotstatusArgs(args []string) error {

	if args[0] != "botstatus" {
		return errors.New("BAD USAGE")

	}

	if args[1] != "on" && args[1] != "off" {
		return errors.New("BAD USAGE")
	}

	return nil
}

func checkCoinsToMonitorArgs(args []string) error {

	if !strings.Contains(args[0], "/") {
		return errors.New("BAD USAGE")
	}

	if args[1] != "on" && args[1] != "off" {
		return errors.New("BAD USAGE")
	}
	number, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.New("BAD USAGE")
	}

	if number < 1 || number > 100 {
		return errors.New("BAD USAGE")
	}

	return nil
}

func getArgsFromText(text string) []string {

	args := strings.Fields(text)
	args = args[1:]

	if len(args) == 2 {
		if err := checkBotstatusArgs(args); err != nil {
			return nil
		}
	}

	if len(args) == 3 {
		if err := checkCoinsToMonitorArgs(args); err != nil {
			return nil
		}
	}

	return args
}

func addToTable(pair, status, percent string) (string, error) {
	request := `
	INSERT INTO coins_to_monitor (pair, status, percent)
	VALUES ($1, $2, $3)
	RETURNING id`
	id := 0

	err := dbClient.QueryRow(request, strings.ToUpper(pair), status, percent).Scan(&id)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return pair + " Added to DB", nil
}

func updateCoin(pair, status, percent string) (string, error) {

	query := `UPDATE coins_to_monitor
	SET pair=$1, status=$2, percent=$3
	WHERE pair = $1
	RETURNING id;`

	id := 0
	err := dbClient.QueryRow(query, strings.ToUpper(pair), status, percent).Scan(&id)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return pair + " updated", nil
}

func updateBotStatus(status string) (string, error) {

	query := `UPDATE lambotrunk_bot_status
	SET bot_status=$1
	WHERE id = 1
	RETURNING bot_status;`

	if status == "on" {
		status = "0"
	} else if status == "off" {
		status = "1"
	}

	botStatus := ""
	err := dbClient.QueryRow(query, status).Scan(&botStatus)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if botStatus == "0" {
		botStatus = "on"
	} else if status == "1" {
		botStatus = "off"
	}

	return "New status set : " + botStatus, nil
}

func getCoinInDB(pair string) (sCoinsToMonitor, error) {

	var coinsToMonitor sCoinsToMonitor

	request := `SELECT * FROM coins_to_monitor where pair='` + strings.ToUpper(pair) + `'`

	row := dbClient.QueryRow(request)
	err := row.Scan(&coinsToMonitor.id, &coinsToMonitor.pair, &coinsToMonitor.status, &coinsToMonitor.percent)
	switch err {
	case sql.ErrNoRows:
		return sCoinsToMonitor{}, nil
	case nil:
		return coinsToMonitor, nil
	default:
		logrus.Error(err)
		return sCoinsToMonitor{}, err
	}
}

func storeCoinInDB(coinToMonitor sCoinsToMonitor, pair, status, percent string) (string, error) {

	if coinToMonitor == (sCoinsToMonitor{}) {
		return addToTable(pair, status, percent)
	} else if coinToMonitor.pair == strings.ToUpper(pair) {
		return updateCoin(pair, status, percent)
	}

	return "", nil
}

func updateCoinsToMonitor(args []string) (string, error) {

	if err := createTable(dbClient, "coins_to_monitor", `id  SERIAL PRIMARY KEY,
		pair TEXT,
		status TEXT,
		percent TEXT`); err != nil {
		return "", err
	}

	pair := args[0]
	status := args[1]
	percent := args[2]

	coinsToMonitor, err := getCoinInDB(pair)
	if err != nil {
		return err.Error(), err
	}
	msg, err := storeCoinInDB(coinsToMonitor, pair, status, percent)
	if err != nil {
		return err.Error(), err
	}

	return msg, nil
}

func LamboTrunk(m *tb.Message) (string, error) {

	var err error

	dbClient, err = initDataBase()
	if err != nil {
		return "", err
	}
	defer dbClient.Close()

	args := getArgsFromText(strings.ToLower(m.Text))
	if args == nil {
		return "Usage : /lamboTrunk On/Off\n/lamboTrunk ETH/BUSD on 100", errors.New("BAD USAGE")
	}

	if args[0] == "botstatus" {
		return updateBotStatus(args[1])
	} else {
		return updateCoinsToMonitor(args)
	}
}
