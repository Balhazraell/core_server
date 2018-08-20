package core

type Chunc struct {
	State       int
	Сoordinates [][]int
}

type Room struct {
	Id      int
	Map     []*Chunc
	Clients []*Client
}
