package core

// Временное решение для того что бы отличать id комнат
var maxID int

type Chunc struct {
	State       int
	Сoordinates [][2]int
}

type Room struct {
	Id      int
	Map     []*Chunc
	Clients []*Client
}

func (room *Room) createMap() {
	var step = 100
	var y = 0

	for i := 0; i < 3; i++ {
		var x = 0
		for j := 0; j < 3; j++ {
			chunc := Chunc{}
			chunc.Сoordinates = append(
				chunc.Сoordinates,
				[2]int{x, y},
				[2]int{x + step, y},
				[2]int{x + step, y + step},
				[2]int{x, y + step},
			)

			x += step

			room.Map = append(room.Map, &chunc)
		}
	}
}

func () Start
