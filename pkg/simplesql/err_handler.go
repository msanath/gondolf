package simplesql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mattn/go-sqlite3"
)

type ErrHandler func(error) error

var (
	ErrInsertConflict = errors.New("insert conflict")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidVersion = errors.New("invalid version or version not provided")
	ErrInternal       = errors.New("internal error")
)

// MySQLErrHandler processes MySQL errors and returns a custom StorageError
func MySQLErrHandler(err error) error {
	if err == nil {
		return nil
	}

	// Check if the error is a SQLite error
	// Record not found
	if err == sql.ErrNoRows {
		return ErrRecordNotFound
	}
	// Check if the error is a MySQL error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			// Duplicate entry (unique constraint violation)
			fallthrough
		case 1451:
			// Foreign key constraint violation (cannot delete/update parent row)
			fallthrough
		case 1452:
			// Foreign key constraint violation (cannot add/update child row)
			return fmt.Errorf("%s: %w", mysqlErr.Error(), ErrInsertConflict)
		default:
			// For all other MySQL errors
			return fmt.Errorf("%s: %w", mysqlErr.Error(), ErrInternal)
		}
	}
	return err
}

// SQLiteErrHandler processes SQLite errors and returns a custom StorageError
func SQLiteErrHandler(err error) error {
	if err == nil {
		return nil
	}

	// Check if the error is a SQLite error
	if err == sql.ErrNoRows {
		// Record not found
		return fmt.Errorf("%s: %w", err.Error(), ErrRecordNotFound)
	}

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.Code {
		case sqlite3.ErrConstraint:
			// Unique constraint violation
			return fmt.Errorf("%s: %w", sqliteErr.Error(), ErrInsertConflict)
		case sqlite3.ErrNotFound:
			// Record not found
			return fmt.Errorf("%s: %w", sqliteErr.Error(), ErrRecordNotFound)
		default:
			// For all other SQLite errors
			return fmt.Errorf("%s: %w", sqliteErr.Error(), ErrInternal)
		}
	}

	return err
}

func defaultErrHandler(err error) error {
	return err
}
