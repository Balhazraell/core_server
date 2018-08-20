package core

// В первой итерации вебсокеты будут передавать сообщения на прямую в gameServer.
// но потом надо продумать другую связь, возможно через балансировщик.

type Client struct {
	Id     int
	RoomId int //Здесь должна быть ссылка на Room
}

type gameServer struct {
}

func Start() {

}
