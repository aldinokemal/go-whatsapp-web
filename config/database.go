package config

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

type TableAccount struct {
	AccID          int    `json:"acc_id"`
	AccPhone       string `json:"acc_phone"`
	AccSessionName string `json:"acc_session_name"`
	AccCreatedAt   string `json:"acc_created_at"`
}

func DBOpen() {
	DB, _ = sql.Open("sqlite3", "./storage/database/go_whatsapp.db")
	createTable(DB)
}

func createTable(db *sql.DB) {
	createAccountTable := `CREATE TABLE IF NOT EXISTS account (
		"acc_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"acc_phone" TEXT,
		"acc_session_name" TEXT,
        "acc_created_at" timestamp default (datetime('now','localtime'))
	  );` // SQL Statement for Create Table

	log.Println("Create account table (if not exist)...")
	statement, err := db.Prepare(createAccountTable) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Table Account created")
}

func (acc TableAccount) InsertAccount() (err error) {
	query := `INSERT INTO account(acc_phone, acc_session_name) VALUES (?, ?)`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccPhone, acc.AccSessionName)
	if err != nil {
		return err
	}
	return nil
}

func (acc TableAccount) FindByPhone() (data TableAccount) {
	query := `SELECT * FROM account WHERE acc_phone = ?`
	row := DB.QueryRow(query, acc.AccPhone) // Prepare SQL Statement
	_ = row.Scan(&data.AccID, &data.AccPhone, &data.AccSessionName, &data.AccCreatedAt)
	return
}

func (acc TableAccount) DelByID() (err error) {
	query := `DELETE FROM account WHERE acc_id = ?`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(query, acc.AccID) // Prepare SQL Statement
	if err != nil {
		return err
	}
	return nil
}
