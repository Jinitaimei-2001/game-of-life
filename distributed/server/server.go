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

/** Super-Secret `reversing a string' method we can't allow clients to see. **/
//
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

//
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
//
//type GameoflifeOperations struct {
//}
//

//func (g *GameoflifeOperations) Worker(req stubs.Request, res *stubs.Response) (err error) {
//	curWorld = req.World
//	curTurn = 0
//	for curTurn < req.Turn {
//		mu.Lock()
//		curTurn++
//		curWorld = calculateNextState(0, req.ImageHeight, 0, req.ImageWidth, curWorld)
//		mu.Unlock()
//	}
//	res.Turn = curTurn
//	res.World = curWorld
//	return
//}
//
//func (g *GameoflifeOperations) AliveCount(req stubs.Request, res *stubs.Response) (err error) {
//	mu.Lock()
//	res.World = curWorld
//	res.Turn = curTurn
//	mu.Unlock()
//	return
//}
//
//func (g *GameoflifeOperations) KeyPress(req stubs.Request, res *stubs.Response) (err error) {
//	res.Turn = curTurn
//	res.World = curWorld
//	//os.Exit(0)
//	return
//}

type Server struct {
}

//	func worker(startY, endY, startX, endX int, previousWorld [][]uint8, outChan chan<- [][]uint8) {
//		outChan <- calculateNextState(startY, endY, startX, endX, previousWorld)
//	}

func calculateNextState(startY, endY, startX, endX int, world [][]byte) [][]byte {
	//newWorld := createNewWorld(p.ImageHeight, p.ImageWidth)
	ImageHeight := endY - startY
	ImageWidth := endX - startX
	newWorld := createNewWorld(ImageHeight, ImageWidth)
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			// Initialise the neighboursAlive to 0.
			neighboursAlive := 0
			for i := -1; i < 2; i++ {
				for j := -1; j < 2; j++ {
					// Mark all of the neighbours excluding the cell itself.
					if i == 0 && j == 0 {
						continue
					}
					//If the cell is on the edge of the diagram, mod it to fix the rule of the game.
					//if world[(y+i+ImageHeight)%ImageHeight][(x+j+ImageWidth)%ImageWidth] != 0 {
					//	neighboursAlive += 1
					if world[(y+i+ImageHeight)%ImageHeight][(x+j+ImageWidth)%ImageWidth] != 0 {
						neighboursAlive += 1
					}
				}
			}

			// When the colour is white, the cell status is alive, parameter is 255.
			// When the colour is black, the cell status is dead, parameter is 0.
			if world[y][x] == 255 {
				//neighboursAlive = neighboursAlive - 1
				// If less than 2 or more than 3 neighbours, live cells dead.
				if (neighboursAlive < 2) || (neighboursAlive > 3) {
					newWorld[y-startY][x] = 0

				} else {
					newWorld[y-startY][x] = 255
				}
			}
			if world[y][x] == 0 {
				// If 3 neighbours alive, dead cells alive.
				if neighboursAlive == 3 {
					newWorld[y-startY][x] = 255
				} else {
					newWorld[y-startY][x] = 0
				}
			}
		}
	}
	return newWorld
}

func (s *Server) Update(req stubs.WorkerRequest, res *stubs.WorkerResponse) (err error) {
	curWorld = req.World
	curTurn = 0
	for curTurn < req.Turn {
		mu.Lock()
		curTurn++
		curWorld = calculateNextState(0, req.ImageHeight, 0, req.ImageWidth, curWorld)
		mu.Unlock()
	}
	res.Turn = curTurn
	fmt.Println(res.Turn)
	res.World = curWorld
	return
}

func main() {
	pAddr := flag.String("port", "8010", "Port to listen on")
	flag.Parse()
	rpc.Register(new(Server))
	listener, err := net.Listen("tcp", ":"+*pAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("listening on %s", listener.Addr().String())
	defer listener.Close()
	rpc.Accept(listener)
}
