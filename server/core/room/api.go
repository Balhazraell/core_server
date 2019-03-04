package room

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"ClientConnect":    APIClientConnect,
	"ClientDisconnect": APIClientDisconnect,
	"SetChunckState":   APISetChunckState,
}

func APIClientConnect(message string) {

}

func APIClientDisconnect(message string) {

}

func APISetChunckState(message string) {

}
