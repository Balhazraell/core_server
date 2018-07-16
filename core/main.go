package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func init() {
	fsJS := http.FileServer(http.Dir("../static/js/"))
	fsCSS := http.FileServer(http.Dir("../static/css/"))

	http.Handle("/js/", http.StripPrefix("/js", fsJS))
	http.Handle("/css/", http.StripPrefix("/css", fsCSS))

	http.HandleFunc("/", returnIndex)
}

func returnIndex(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("../static/index.html")

	if err != nil {
		fmt.Fprintf(response, err.Error())
	}

	templateErr := t.ExecuteTemplate(response, "index.html", nil)

	if templateErr != nil {
		fmt.Fprintf(response, templateErr.Error())
		fmt.Fprintf(response, t.DefinedTemplates())
	}
}

func main() {
	http.ListenAndServe(":8081", nil)
	fmt.Println("Server is started...")
}
