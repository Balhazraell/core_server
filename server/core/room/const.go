package room

// Состояния присущие чанкам комнаты.
// ChuncStateEmpty = 0 - Пустая ячейка.
// ChuncStateCross = 1 - Крестик.
//ChuncStateZero = 2 - Нолик.
const (
	ChuncStateEmpty = iota
	ChuncStateCross
	ChuncStateZero
)
