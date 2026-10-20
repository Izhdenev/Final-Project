package db

import (
	"database/sql"
	"fmt"
	"time"
)

const DateLayout = "20060102"

const searchDateLayout = "02.01.2006"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка вставки задачи: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID: %w", err)
	}
	return id, nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	baseQuery := `SELECT id, date, title, comment, repeat FROM scheduler`

	var (
		rows *sql.Rows
		err  error
	)

	switch {
	case search == "":
		rows, err = db.Query(baseQuery+` ORDER BY date LIMIT ?`, limit)

	default:
		if t, e := time.Parse(searchDateLayout, search); e == nil {
			rows, err = db.Query(
				baseQuery+` WHERE date = ? ORDER BY date LIMIT ?`,
				t.Format(DateLayout), limit,
			)
		} else {
			like := "%" + search + "%"
			rows, err = db.Query(
				baseQuery+` WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
				like, like, limit,
			)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка выборки задач: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка чтения задачи: %w", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка обхода строк: %w", err)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var (
		t     Task
		intID int64
	)
	err := db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&intID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	t.ID = fmt.Sprintf("%d", intID)
	return &t, nil
}

func UpdateTask(task *Task) error {
	res, err := db.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновлённых строк: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удалённых строк: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := db.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновлённых строк: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
