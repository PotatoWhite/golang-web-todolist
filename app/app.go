package app

import (
	"github.com/gorilla/mux"
	"github.com/potatowhite/web/golang-todolist/model"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
)

var rd *render.Render

func MakeHandler() http.Handler {
	//todoMap = make(map[int]*model.Todo)
	rd = render.New()

	//addTestTodos()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")
	return router
}

func completeTodoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])
	complete := request.FormValue("complete") == "true"
	if ok := model.CompleteTodo(id, complete); ok {
		rd.JSON(writer, http.StatusOK, Success{Success: true})
		return
	}

	rd.JSON(writer, http.StatusInternalServerError, Success{Success: false})
}

type Success struct {
	Success bool `json:"success"`
}

func removeTodoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	if ok := model.RemoveTodo(id); ok {
		rd.JSON(writer, http.StatusNoContent, Success{Success: true})
		return
	} else {
		rd.JSON(writer, http.StatusNoContent, Success{Success: false})
	}
}

func addTodoHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	rd.JSON(writer, http.StatusCreated, model.AddTodo(name))
}

func getTodoListHandler(writer http.ResponseWriter, request *http.Request) {
	rd.JSON(writer, http.StatusOK, model.GetTodos())
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/todo.html", http.StatusTemporaryRedirect)
}
