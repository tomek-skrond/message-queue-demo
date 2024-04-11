package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", dashboardHandler)

	log.Fatalln(http.ListenAndServe(":8080", r))
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalln(err)
		return
	}
	tmpl.Execute(w, "")
}
