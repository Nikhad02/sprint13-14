package db

import (
	"database/sql"
	"errors"
	"fmt"
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

func GetTask(id string) (*Task, error) {
	// Переводим строковый ID в число для базы данных
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("неверный формат идентификатора")
	}

	var t Task
	var dbID int

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	// QueryRow выполняет запрос и сразу готовит данные для Scan
	err = DB.QueryRow(query, idInt).Scan(&dbID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}

	// Записываем строковый ID обратно в структуру
	t.ID = strconv.Itoa(dbID)
	return &t, nil
}

func UpdateTask(task *Task) error {
	idInt, err := strconv.Atoi(task.ID)
	if err != nil {
		return errors.New("неверный формат идентификатора")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, idInt)
	if err != nil {
		return err
	}

	// Проверяем, изменилось ли что-то в базе
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("неверный формат идентификатора")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, idInt)
	if err != nil {
		return err
	}

	// Проверяем, было ли вообще что удалять
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateTaskDate(id string, nextDate string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("неверный формат идентификатора")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, nextDate, idInt)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
