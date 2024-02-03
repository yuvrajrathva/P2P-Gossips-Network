package main

import (
	"github.com/yuvrajrathva/P2P-Gossips-Network/controllers"
)

func main(){
	controllers.CreateSeed()
	controllers.CreatePeer("127.0.0.1", 8080)
}