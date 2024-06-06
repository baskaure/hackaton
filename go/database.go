package hackaton

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USER (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        pseudo TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS HISTORY (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        latitude REAL NOT NULL,
        longitude REAL NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(user_id) REFERENCES USER(id)
    )`)
    if err != nil {
        log.Fatal(err)
    }
	return db
}

func InsertUser(db *sql.DB, username, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    stmt, err := db.Prepare("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    _, err = stmt.Exec(username, email, hashedPassword)
    return err
}

func AuthenticateUser(db *sql.DB, username, password string) (bool, error) {
    var hashedPassword string
    err := db.QueryRow("SELECT password FROM USER WHERE pseudo = ?", username).Scan(&hashedPassword)
    if err != nil {
        return false, err
    }
    err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil, err
}

func InsertHistory(db *sql.DB, userId int, latitude, longitude float64) error {
	stmt, err := db.Prepare("INSERT INTO HISTORY (user_id, latitude, longitude) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, latitude, longitude)
	return err
}