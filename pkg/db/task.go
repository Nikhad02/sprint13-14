package db

import (
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	// Получаем ID, который SQLite сгенерировал автоматически (AUTOINCREMENT)
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int, search string) ([]Task, error) {
	tasks := make([]Task, 0)

	var query string
	var args []interface{}

	if search == "" {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		args = append(args, limit)
	} else {
		parsedTime, err := time.Parse("02.01.2006", search)
		if err == nil {
			targetDate := parsedTime.Format("20060102")
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
			args = append(args, targetDate, limit)
		} else {
			likeSearch := "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
			args = append(args, likeSearch, likeSearch, limit)
		}
	}
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idInt int
		var t Task

		// Сканируем данные из БД. ID забираем как число.
		err = rows.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		// Переводим числовой ID в строку для соответствия структуре Task
		t.ID = strconv.Itoa(idInt)

		tasks = append(tasks, t)
	}

	// Проверяем, не было ли ошибок во время итерации по строкам
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
