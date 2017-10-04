package scripting

import (
	"database/sql"
	"time"
)

// AddTag takes a name and the contents of a script and stores it in a database.
func AddTag(name string, script string, author string, Database *sql.DB) error {
	transaction, err := Database.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare(`
		INSERT INTO scripts(name, content, author, date, uses) VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(name, script, author, time.Now(), 0)
	if err != nil {
		return err
	}

	transaction.Commit()
	return nil
}

// DeleteTag takes a name and deletes that script from the database.
func DeleteTag(name string, Database *sql.DB) error {
	statement, err := Database.Prepare("DELETE FROM scripts WHERE name = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

// GetRawScript returns the full contents of the script
func GetRawScript(name string, Database *sql.DB) (string, error) {
	statement, err := Database.Prepare("SELECT content FROM scripts WHERE name = ?")
	if err != nil {
		return "", err
	}
	content := ""
	err = statement.QueryRow(name).Scan(&content)
	if err != nil {
		return "", err
	}

	return content, nil
}
