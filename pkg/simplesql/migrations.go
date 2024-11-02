package simplesql

import (
	"sort"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Version int
	Up      string
	Down    string
}

func getCurrentSchemaVersion(db *sqlx.DB) (int, error) {
	var version int
	err := db.Get(&version, "SELECT version FROM schema_version")
	return version, err
}

func setCurrentSchemaVersion(db *sqlx.DB, version int) error {
	_, err := db.Exec("DELETE FROM schema_version")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO schema_version (version) VALUES (?)", version)
	return err
}

func (d *Database) ApplyMigrations(schemaMigrations []Migration) error {
	_, err := d.DB.Exec("CREATE TABLE IF NOT EXISTS schema_version (version INTEGER PRIMARY KEY);")
	if err != nil {
		return err
	}

	currentVersion, err := getCurrentSchemaVersion(d.DB)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			currentVersion = 0
		} else {
			return err
		}
	}

	// Sort the list of schemas by Version
	sort.Slice(schemaMigrations, func(i, j int) bool {
		return schemaMigrations[i].Version < schemaMigrations[j].Version // Ascending order
	})

	for _, s := range schemaMigrations {
		if s.Version > currentVersion {
			_, err := d.DB.Exec(s.Up)
			if err != nil {
				return err
			}

			if err := setCurrentSchemaVersion(d.DB, s.Version); err != nil {
				return err
			}
		}
	}

	return nil
}
