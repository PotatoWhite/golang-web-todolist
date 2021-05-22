package main

import (
	"github.com/potatowhite/web/golang-todolist/app"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	router := app.MakeHandler()
	neg := negroni.Classic()
	neg.UseHandler(router)

	err := http.ListenAndServe(":3000", neg)
	if err != nil {
		panic(err)
	}
}
