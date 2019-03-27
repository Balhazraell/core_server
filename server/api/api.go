package api

var API websocketsCoreAPI

func init() {
	API = websocketsCoreAPI{
		// core
		make(chan int),
		make(chan int),
		make(chan SetChunckStateStruct),
		make(chan ChangeRoomStructure),

		// websockets
		make(chan UpdateClientsMapStruct),
		make(chan NewClientIsConnectedStruct),
		make(chan SendErrorToСlientStruct),
	}
}

type websocketsCoreAPI struct {
	// core
	ClientConnectionChl chan int
	ClientDisconnectChl chan int
	SetChunckStateChl   chan SetChunckStateStruct
	ChangeRoomChl       chan ChangeRoomStructure

	//websockets
	UpdateClientsMapChl     chan UpdateClientsMapStruct
	NewClientIsConnectedChl chan NewClientIsConnectedStruct
	SendErrorToСlientChl    chan SendErrorToСlientStruct
}

// core
type SetChunckStateStruct struct {
	ClientID int
	ChuncID  int
}

type ChangeRoomStructure struct {
	ClientID int
	RoomID   int
}

// websockets
type UpdateClientsMapStruct struct {
	GameMap    []byte
	ClientsIDs []int
}

type RoomData struct {
	ID   int
	Name string
}

type NewClientIsConnectedStruct struct {
	ClientID     int
	RoomsCatalog []RoomData
}

type SendErrorToСlientStruct struct {
	ClientID int
	Message  string
}
