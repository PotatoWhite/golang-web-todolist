package app

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/potatowhite/web/golang-todolist/app/oauth/google"
	"github.com/potatowhite/web/golang-todolist/model"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var rd = render.New()

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

func (self *AppHandler) Close() {
	self.db.Close()
}

func MakeHandler(filePath string) *AppHandler {
	router := mux.NewRouter()
	neg := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(checkSignIn), negroni.NewStatic(http.Dir("public")))
	neg.UseHandler(router)

	app := &AppHandler{
		Handler: neg,
		db:      model.NewDBHandler(filePath),
	}

	// register session store
	google.SetSessionHandler(addSessionId)

	router.HandleFunc("/", app.indexHandler)
	router.HandleFunc("/todos", app.getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", app.addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", app.removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", app.completeTodoHandler).Methods("GET")

	//oauth
	router.HandleFunc("/auth/google/login", google.RedirectToGoogleLoginPage)
	router.HandleFunc("/auth/google/callback", google.CallBackOAuthResultAndPrintUserInfo)

	return app
}

func checkSignIn(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
	if strings.Contains(request.URL.Path, "/signin") || strings.Contains(request.URL.Path, "/auth/") {
		next(writer, request)
		return
	}

	_, err := getSessionId(request)
	if err == nil {
		next(writer, request)
		return
	}

	http.Redirect(writer, request, "/signin.html", http.StatusTemporaryRedirect)
}

func (self *AppHandler) completeTodoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])
	complete := request.FormValue("complete") == "true"
	if ok := self.db.CompleteTodo(id, complete); ok {
		rd.JSON(writer, http.StatusOK, Success{Success: true})
		return
	}

	rd.JSON(writer, http.StatusInternalServerError, Success{Success: false})
}

type Success struct {
	Success bool `json:"success"`
}

func (self *AppHandler) removeTodoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	if ok := self.db.RemoveTodo(id); ok {
		rd.JSON(writer, http.StatusNoContent, Success{Success: true})
		return
	} else {
		rd.JSON(writer, http.StatusNoContent, Success{Success: false})
	}
}

func (self *AppHandler) addTodoHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	rd.JSON(writer, http.StatusCreated, self.db.AddTodo(name))
}

func (self *AppHandler) getTodoListHandler(writer http.ResponseWriter, request *http.Request) {
	rd.JSON(writer, http.StatusOK, self.db.GetTodos())
}

func (self *AppHandler) indexHandler(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/todo.html", http.StatusTemporaryRedirect)
}

func getSessionId(request *http.Request) (interface{}, error) {
	session, err := store.Get(request, "session")
	if err != nil {
		return "", err
	}

	id := session.Values["id"]
	if id == nil {
		return "", errors.New("Empty session")
	}

	return id, nil
}

func addSessionId(writer http.ResponseWriter, request *http.Request, id string) error {
	// get session info
	session, err := store.Get(request, "session")
	if err != nil {
		return err
	}

	// store session id
	session.Values["id"] = id
	err = session.Save(request, writer)
	if err != nil {
		return err
	}

	return nil
}
