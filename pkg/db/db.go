package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// schema содержит SQL для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init открывает БД и создаёт таблицу/индекс, если файла не существовало
func Init(dbFile string) error {
	// Проверяем существование файла
	_, err := os.Stat(dbFile)
	install := err != nil

	// Открываем БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файла не было — выполняем создание таблицы и индекса
	if install {
		log.Println("Creating database schema...")
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	log.Println("Database initialized:", dbFile)
	return nil
}
