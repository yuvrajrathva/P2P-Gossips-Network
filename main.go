package main

import (
	"github.com/yuvrajrathva/P2P-Gossips-Network/controllers"
)

func main(){
	controllers.CreateSeed()
	go controllers.CreatePeer("127.0.0.1", 8000)
	go controllers.CreatePeer("127.0.0.1", 8001)
	go controllers.CreatePeer("127.0.0.1", 8002)
}