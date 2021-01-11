package config

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB

type TableAccount struct {
	AccID          int            `json:"acc_id"`
	AccAppID       string         `json:"acc_app_id"`
	AccWaID        sql.NullString `json:"acc_wa_id"`
	AccQrName      sql.NullString `json:"acc_qr_name"`
	AccSessionName sql.NullString `json:"acc_session_name"`
	AccCreatedAt   sql.NullTime   `json:"acc_created_at"`
}

func DBOpen() {
	DB, _ = sql.Open("sqlite3", "./storage/database/go_whatsapp.db")
	createTable(DB)
}

func createTable(db *sql.DB) {
	createAccountTable := `CREATE TABLE IF NOT EXISTS account (
		"acc_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"acc_app_id" TEXT UNIQUE,
		"acc_wa_id" TEXT,
		"acc_qr_name" TEXT,
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
	query := `INSERT INTO account(acc_app_id, acc_qr_name,acc_session_name) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccAppID, acc.AccQrName, acc.AccSessionName)
	if err != nil {
		return err
	}
	return nil
}

func (acc TableAccount) FindAll() (data []TableAccount) {
	query := `SELECT * FROM account WHERE acc_session_name IS NOT NULL`
	rows, _ := DB.Query(query)
	defer rows.Close()

	for rows.Next() {
		var x TableAccount
		err := rows.Scan(&x.AccID, &x.AccAppID, &x.AccWaID, &x.AccQrName, &x.AccSessionName, &x.AccCreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, x)
	}
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (acc TableAccount) FindByAppID() (data TableAccount) {
	query := `SELECT * FROM account WHERE acc_app_id = ?`
	row := DB.QueryRow(query, acc.AccAppID)
	fmt.Print()
	err := row.Scan(&data.AccID, &data.AccAppID, &data.AccWaID, &data.AccQrName, &data.AccSessionName, &data.AccCreatedAt)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func (acc TableAccount) UpdateSessionByAppID() (err error) {
	query := `UPDATE account SET acc_session_name = ?, acc_wa_id = ? WHERE acc_app_id = ?`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccSessionName, acc.AccWaID.String, acc.AccAppID) // Prepare SQL Statement
	if err != nil {
		return err
	}
	return nil
}

func (acc TableAccount) DelByAppID() (err error) {
	query := `DELETE FROM account WHERE acc_app_id = ?`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccAppID) // Prepare SQL Statement
	if err != nil {
		return err
	}
	return nil
}
