package app

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

var rd *render.Render

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var todoMap map[int]*Todo
var lastIndex uint = 0

func MakeHandler() http.Handler {
	todoMap = make(map[int]*Todo)
	rd = render.New()

	addTestTodos()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", addTodoHandler).Methods("POST")
	return router
}

func addTodoHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	id := len(todoMap) + 1
	todo := &Todo{Id: id, Name: name, CreatedAt: time.Now()}
	todoMap[id] = todo

	rd.JSON(writer, http.StatusCreated, todo)

}

func addTestTodos() {
	todoMap[1] = &Todo{1, "Buy a milk", false, time.Now()}
	todoMap[2] = &Todo{2, "Exercise", true, time.Now()}
	todoMap[3] = &Todo{3, "Homework", false, time.Now()}

}

func getTodoListHandler(writer http.ResponseWriter, request *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}

	rd.JSON(writer, http.StatusOK, list)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/todo.html", http.StatusTemporaryRedirect)
}
