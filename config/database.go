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
	AccPhone       string         `json:"acc_phone"`
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
		"acc_phone" TEXT UNIQUE,
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
	query := `INSERT INTO account(acc_phone, acc_qr_name,acc_session_name) VALUES (?, ?, ?)`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccPhone, acc.AccQrName, acc.AccSessionName)
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
		err := rows.Scan(&x.AccID, &x.AccPhone, &x.AccQrName, &x.AccSessionName, &x.AccCreatedAt)
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

func (acc TableAccount) FindByPhone() (data TableAccount) {
	query := `SELECT * FROM account WHERE acc_phone = ?`
	row := DB.QueryRow(query, acc.AccPhone)
	fmt.Print()
	err := row.Scan(&data.AccID, &data.AccPhone, &data.AccQrName, &data.AccSessionName, &data.AccCreatedAt)
	if err != nil {
		fmt.Println(err.Error())
	}

	return data
}

func (acc TableAccount) UpdateSessionByPhone() (err error) {
	query := `UPDATE account SET acc_session_name = ? WHERE acc_phone = ?`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccSessionName, acc.AccPhone) // Prepare SQL Statement
	if err != nil {
		return err
	}
	return nil
}

func (acc TableAccount) DelByPhone() (err error) {
	query := `DELETE FROM account WHERE acc_phone = ?`
	statement, err := DB.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec(acc.AccPhone) // Prepare SQL Statement
	if err != nil {
		return err
	}
	return nil
}
