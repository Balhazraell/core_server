package api

var API websocketsCoreAPI

func init() {
	API = websocketsCoreAPI{
		// core
		make(chan int),
		make(chan int),
		make(chan SetChunckStateStruct),

		// websockets
		make(chan UpdateClientsMapStruct),
		make(chan NewClientIsConnectedStruct),
	}
}

type websocketsCoreAPI struct {
	// core
	ClientConnectionChl chan int
	ClientDisconnectChl chan int
	SetChunckStateChl   chan SetChunckStateStruct

	//websockets
	UpdateClientsMapChl     chan UpdateClientsMapStruct
	NewClientIsConnectedChl chan NewClientIsConnectedStruct
}

// core
type SetChunckStateStruct struct {
	ClientID int
	ChuncID  int
}

// websockets
type UpdateClientsMapStruct struct {
	GameMap    []byte
	ClientsIDs []int
}

type NewClientIsConnectedStruct struct {
	ClientID  int
	ClientMap []byte
}
