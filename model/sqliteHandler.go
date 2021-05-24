package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteHandler struct {
	db *sql.DB
}

func (self *sqliteHandler) close() {
	defer self.db.Close()
}

func (self *sqliteHandler) getTodos() []*Todo {
	return nil
}
func (self *sqliteHandler) addTodo(name string) *Todo {
	return nil
}
func (self *sqliteHandler) removeTodo(id int) bool {
	return false
}
func (self *sqliteHandler) completeTodo(id int, complete bool) bool {
	return false
}

func newSqliteHandler() dbHandler {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			name TEXT,
    			completed BOOLEAN,
    			createAt DATETIME)
    			`)
	if err != nil {
		panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}

	return &sqliteHandler{db: db}
}
