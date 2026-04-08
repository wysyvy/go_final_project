package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

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

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Init открывает БД и создаёт таблицу/индекс, если файла не существовало
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := err != nil

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

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

// AddTask добавляет задачу в БД и возвращает id
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetTask возвращает задачу по id
func GetTask(id string) (*Task, error) {
	var task Task
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask удаляет задачу по id
func DeleteTask(id string) error {
	result, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateTaskDate обновляет только дату
func UpdateTaskDate(id string, date string) error {
	result, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
