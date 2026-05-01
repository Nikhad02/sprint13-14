package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // Импорт драйвера без прямого использования
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

// Init открывает базу данных и создает таблицу, если файла не было
func Init(dbFile string) error {
	var install bool

	// Проверяем существование файла
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	// Открываем базу данных
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем соединение
	if err = DB.Ping(); err != nil {
		return err
	}

	// Если файла не было, выполняем установку схемы
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
