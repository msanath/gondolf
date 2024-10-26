package test

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// NewTestSQLiteDB creates a new SQLite database for testing.
// It is an sqlite database in memory.
func NewTestSQLiteDB() (*sqlx.DB, error) {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign key constraints: %w", err)
	}
	return db, nil
}

// NewTestMySQLDB creates a new MySQL database for testing.
// It neeeds mysql to be running on localhost:3306 with root user having no password.
func NewTestMySQLDB() (*sqlx.DB, error) {
	// local mysql with no password
	dsn := "root:@tcp(127.0.0.1:3306)/?parseTime=true"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Step 2: Delete the database if it already exists
	dbName := "test_db"
	dropDBQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	_, err = db.Exec(dropDBQuery)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to drop database: %w", err)
	}

	// Force MySQL to flush and ensure all commands are committed
	_, err = db.Exec("FLUSH TABLES")
	if err != nil {
		return nil, fmt.Errorf("failed to flush tables: %w", err)
	}

	// Step 3: Create the database
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Step 3: Now connect to the newly created database
	// This is important: You must connect to the database in which you want to create tables.
	dsnWithDB := fmt.Sprintf("root:@tcp(127.0.0.1:3306)/%s", dbName)
	dbWithDB, err := sqlx.Connect("mysql", dsnWithDB)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return dbWithDB, nil
}
