package model

import "time"

type dbHandler interface {
	getTodos() []*Todo
	addTodo(name string) *Todo
	removeTodo(id int) bool
	completeTodo(id int, complete bool) bool
}

type MemoryHandler struct {
	todoMap   map[int]*Todo
	lastIndex int
}

func (self *MemoryHandler) getTodos() []*Todo {
	list := []*Todo{}
	for _, v := range self.todoMap {
		list = append(list, v)
	}
	return list
}

func (self *MemoryHandler) addTodo(name string) *Todo {
	self.lastIndex++
	self.todoMap[self.lastIndex] = &Todo{Id: self.lastIndex, Name: name, CreatedAt: time.Now(), Completed: false}
	return self.todoMap[self.lastIndex]
}

func (self *MemoryHandler) removeTodo(id int) bool {
	if _, ok := self.todoMap[id]; ok {
		delete(self.todoMap, id)
		return true
	}

	return false
}

func (self *MemoryHandler) completeTodo(id int, complete bool) bool {
	if todo, ok := self.todoMap[id]; ok {
		todo.Completed = complete
		return true
	}

	return false
}

func newMemoryHandler() dbHandler {
	m := &MemoryHandler{}
	m.todoMap = make(map[int]*Todo)
	m.lastIndex = 0

	return m
}
