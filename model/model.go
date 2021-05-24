package model

import (
	"time"
)

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var handler dbHandler

func init() {
	//handler = newMemoryHandler()
	handler = newSqliteHandler()
}

func GetTodos() []*Todo {
	return handler.getTodos()
}

func AddTodo(name string) *Todo {
	return handler.addTodo(name)
}

func RemoveTodo(id int) bool {
	return handler.removeTodo(id)
}

func CompleteTodo(id int, complete bool) bool {
	return handler.completeTodo(id, complete)
}
