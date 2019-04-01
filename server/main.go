package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./core"
	"./logger"
	"./websockets"
)

func init() {
	fsJS := http.FileServer(http.Dir("../static/js/src/dist/"))
	fsCSS := http.FileServer(http.Dir("../static/css/"))

	http.Handle("/js/", http.StripPrefix("/js", fsJS))
	http.Handle("/css/", http.StripPrefix("/css", fsCSS))

	http.HandleFunc("/", returnIndex)
}

func returnIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if r.URL.Path == "/" {
		t, err := template.ParseFiles("../static/index.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		templateErr := t.ExecuteTemplate(w, "index.html", nil)

		if templateErr != nil {
			fmt.Fprintf(w, templateErr.Error())
			fmt.Fprintf(w, t.DefinedTemplates())
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>We could not find the page you "+
			"were looking for :(</h1>"+
			"<p>Please email us if you keep being sent to an "+
			"invalid page.</p>")
	}

}

func main() {
	defer logger.InfoPrint("Сервер закончил работу.")
	// Запускаем логгер.
	ok := logger.InitLogger()

	if ok {
		fmt.Println("Logger - YES")
	} else {
		fmt.Println("Logger - NO")
	}

	// Запускаем игровой сервер.
	core.ServerStart()

	// Стартуем сервер websocket.
	websockets.Start()

	// Стартуем сервер статики. Стартуем его последним.
	// ListenAndServe - ждем завершения, по этому код дальше не выполняется.
	logger.InfoPrint("Сервер запущен и готов к работе.")
	http.ListenAndServe(":8081", nil)
}
