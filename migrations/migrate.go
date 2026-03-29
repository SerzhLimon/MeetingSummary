package migrations

import (
	"database/sql"
	"fmt"
)

func Up(db *sql.DB) error {
    if db == nil {
        return fmt.Errorf("db not init")
    }

    tx, err := db.Begin()
    if err != nil {
        return err
    }

    // Проверяем, существует ли тип 'status'
    var typeExists bool
    err = tx.QueryRow(`
        SELECT EXISTS (
            SELECT 1
            FROM pg_type
            WHERE typname = 'status'
        )
    `).Scan(&typeExists)

    if err != nil {
        tx.Rollback()
        return err
    }

    if !typeExists {
        if _, err := tx.Exec(createEnum); err != nil {
            tx.Rollback()
            return err
        }
    }

    if _, err := tx.Exec(createTable); err != nil {
        tx.Rollback()
        return err
    }

    if err = tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }

    return nil
}

func Down(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db not init")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(dropTable); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(dropEnum); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
