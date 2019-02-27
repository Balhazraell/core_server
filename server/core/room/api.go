package room

// APIMetods - Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"ClientConnect": APIClientConnect,
}

// APIClientConnect - Подключение нового пользователя.
func APIClientConnect(message string) {

}
