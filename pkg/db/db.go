package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Init(dbFile string) error {
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	return nil
}

// GetAllTasks возвращает список задач
func GetAllTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTaskByID возвращает задачу по id
func GetTaskByID(id string) (*Task, error) {
	var task Task
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// AddTask добавляет задачу
func AddTask(task *Task) (int64, error) {
	result, err := DB.Exec(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	result, err := DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// DeleteTaskByID удаляет задачу по ID
func DeleteTaskByID(id string) error {
	result, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// UpdateTaskDate обновляет только дату задачи
func UpdateTaskDate(id string, date string) error {
	result, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
