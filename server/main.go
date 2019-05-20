package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./core"
	"./websockets"
	"github.com/Balhazraell/logger"
)

func init() {
	fsJS := http.FileServer(http.Dir("../static/js/src/dist/"))
	fsCSS := http.FileServer(http.Dir("../static/css/"))

	http.Handle("/js/", http.StripPrefix("/js", fsJS))
	http.Handle("/css/", http.StripPrefix("/css", fsCSS))

	http.HandleFunc("/", returnIndex)
}

type view struct {
	temp *template.Template
}

func createView(paths ...string) *view {
	//!TODO: Необходимо реорганизовать папки.
	paths = append(paths,
		"../static/views/layouts/head.gohtml",
		"../static/views/layouts/main.gohtml",
		"../static/views/layouts/footer.gohtml",
	)

	t, err := template.ParseFiles(paths...)
	if err != nil {
		logger.ErrorPrintf("Не смогли создать страницу %v", err)
	}

	return &view{temp: t}
}

func returnIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		gameWindow := createView("../static/views/game_window.gohtml")
		err := gameWindow.temp.ExecuteTemplate(w, "main", nil)
		if err != nil {
			logger.InfoPrintf("При склейке шаблонов произошла ошибка %v:", err)
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
