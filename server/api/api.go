package api

var API websocketsCoreAPI

func init() {
	API = websocketsCoreAPI{
		// core
		ClientConnectionChl: make(chan int),
		ClientDisconnectChl: make(chan int),
		SetChunckStateChl:   make(chan SetChunckStateStruct),
		ChangeRoomChl:       make(chan ChangeRoomStructure),

		// websockets
		UpdateClientsMapChl:     make(chan UpdateClientsMapStruct),
		NewClientIsConnectedChl: make(chan NewClientIsConnectedStruct),
		SendErrorToСlientChl:    make(chan SendErrorToСlientStruct),
		UpdateRoomsCatalog:      make(chan RoomsCatalogStruct),
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
	UpdateRoomsCatalog      chan RoomsCatalogStruct
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

type RoomsCatalogStruct struct {
	ClientIDs    []int
	RoomsCatalog []RoomData
}
