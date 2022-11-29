package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"uk.ac.bris.cs/gameoflife/stubs"
)

var curWorld [][]byte
var curTurn int
var mu sync.Mutex

func createNewWorld(height, width int) [][]byte {
	World := make([][]byte, height)
	for v := range World {
		World[v] = make([]byte, width)
	}
	return World
}

//func calculateNextState(startY, endY, startX, endX int, world [][]byte) [][]byte {
//	//newWorld := createNewWorld(p.ImageHeight, p.ImageWidth)
//	ImageHeight := endY - startY
//	ImageWidth := endX - startX
//	newWorld := createNewWorld(ImageHeight, ImageWidth)
//	for y := startY; y < endY; y++ {
//		for x := startX; x < endX; x++ {
//			// Initialise the neighboursAlive to 0.
//			neighboursAlive := 0
//			for i := -1; i < 2; i++ {
//				for j := -1; j < 2; j++ {
//					// Mark all of the neighbours excluding the cell itself.
//					if i == 0 && j == 0 {
//						continue
//					}
//					//If the cell is on the edge of the diagram, mod it to fix the rule of the game.
//					//if world[(y+i+ImageHeight)%ImageHeight][(x+j+ImageWidth)%ImageWidth] != 0 {
//					//	neighboursAlive += 1
//					if world[(y+i+ImageHeight)%ImageHeight][(x+j+ImageWidth)%ImageWidth] != 0 {
//						neighboursAlive += 1
//					}
//				}
//			}
//
//			// When the colour is white, the cell status is alive, parameter is 255.
//			// When the colour is black, the cell status is dead, parameter is 0.
//			if world[y][x] == 255 {
//				//neighboursAlive = neighboursAlive - 1
//				// If less than 2 or more than 3 neighbours, live cells dead.
//				if (neighboursAlive < 2) || (neighboursAlive > 3) {
//					newWorld[y-startY][x] = 0
//
//				} else {
//					newWorld[y-startY][x] = 255
//				}
//			}
//			if world[y][x] == 0 {
//				// If 3 neighbours alive, dead cells alive.
//				if neighboursAlive == 3 {
//					newWorld[y-startY][x] = 255
//				} else {
//					newWorld[y-startY][x] = 0
//				}
//			}
//		}
//	}
//	return newWorld
//}

type Broker struct {
}

func (b *Broker) AliveCount(req stubs.Request, res *stubs.Response) (err error) {
	mu.Lock()
	res.World = curWorld
	res.Turn = curTurn
	mu.Unlock()
	return
}

func (b *Broker) KeyPress(req stubs.Request, res *stubs.Response) (err error) {
	res.Turn = curTurn
	res.World = curWorld
	return
}

func (b *Broker) MakeCallWorker(req stubs.Request, res *stubs.Response) (err error) {
	client, err := rpc.Dial("tcp", "127.0.0.1:8010")
	if err != nil {
		return
	}
	defer client.Close()
	request := stubs.WorkerRequest{World: req.World,
		ImageHeight: req.ImageHeight,
		ImageWidth:  req.ImageWidth,
		Turn:        req.Turn,
		Threads:     req.Threads,
	}
	response := stubs.WorkerResponse{}
	err = client.Call(stubs.Worker, request, &response)
	curWorld = response.World
	curTurn = response.Turn
	res.World = curWorld
	res.Turn = curTurn
	fmt.Println(res.Turn)
	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rpc.Register(&Broker{})
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("listening on %s", listener.Addr().String())
	defer listener.Close()
	rpc.Accept(listener)
}
