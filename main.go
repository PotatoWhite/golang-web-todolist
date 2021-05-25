package main

import (
	"github.com/potatowhite/web/golang-todolist/app"
	"log"
	"net/http"
)

func main() {
	app := app.MakeHandler("./todos.db")
	defer app.Close()

	log.Println("Application started")
	err := http.ListenAndServe(":3000", app.Handler)
	if err != nil {
		panic(err)
	}
}
