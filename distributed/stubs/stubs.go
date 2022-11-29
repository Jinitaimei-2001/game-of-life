package stubs

//distributor call broker

//var Gameoflife = "Broker.Update"

var Gameoflife = "Broker.MakeCallWorker"
var KeyPress = "Broker.KeyPress"
var AliveCount = "Broker.AliveCount"

// Broker calls server

var Worker = "Server.Update"

type Request struct {
	World       [][]byte
	Turn        int
	ImageHeight int
	ImageWidth  int
	Threads     int
}

type Response struct {
	World [][]byte
	Turn  int
}

type WorkerResponse struct {
	World [][]byte
	Turn  int
}

type WorkerRequest struct {
	World       [][]byte
	Turn        int
	ImageHeight int
	ImageWidth  int
	Threads     int
}
