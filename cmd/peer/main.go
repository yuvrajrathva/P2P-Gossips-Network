package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/yuvrajrathva/P2P-Gossips-Network/controllers"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Please provide a valid command")
		return
	}

	ip := args[1]
	port, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Please provide a valid port number")
		return
	}
	controllers.CreatePeer(ip, port)
}
