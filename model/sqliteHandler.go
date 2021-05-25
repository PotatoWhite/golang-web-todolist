package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type sqliteHandler struct {
	db *sql.DB
}

func (self *sqliteHandler) Close() {
	defer self.db.Close()
}

func (self *sqliteHandler) GetTodos() []*Todo {
	todos := []*Todo{}
	records, err := self.db.Query("SELECT id, name, completed, createdAt FROM todos")
	if err != nil {
		panic(err)
	}
	defer records.Close()

	for records.Next() {
		var todo Todo
		records.Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}

	return todos
}
func (self *sqliteHandler) AddTodo(name string) *Todo {
	stmt, err := self.db.Prepare("INSERT INTO todos (name, completed, createdAt) VALUES (?,?,datetime('now'))")
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec(name, false)
	if err != nil {
		panic(err)
	}

	id, _ := result.LastInsertId()

	var todo Todo
	todo.Id = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo

}
func (self *sqliteHandler) RemoveTodo(id int) bool {
	stmt, err := self.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}
func (self *sqliteHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := self.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		panic(err)
	}
	result, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}

	rows, _ := result.RowsAffected()
	return rows > 0
}

func newSqliteHandler(dbpath string) DBHandler {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS todos (
    			id INTEGER PRIMARY KEY AUTOINCREMENT,
    			name TEXT,
    			completed BOOLEAN,
    			createdAt DATETIME)
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
