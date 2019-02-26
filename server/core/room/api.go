package room

// Перечень доступных API методов.
var APIMetods = map[string]func(string){
	"roomConnect": RoomConnect,
}
