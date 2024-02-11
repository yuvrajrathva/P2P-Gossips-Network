package controllers

import (
	"fmt"
	"net"
	"sync"
)

func TCPClient(ip string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error connecting to peer - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to peer - IP: %s, Port: %d\n", ip, port)
}

func TCPServer(ip string, port int, wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error starting server - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer listener.Close() 
	defer wg.Done()

	fmt.Printf("Server started - IP: %s, Port: %d\n", ip, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}

		defer conn.Close()
		fmt.Printf("Connection accepted from - IP: %s, Port: %d\n", conn.RemoteAddr().(*net.TCPAddr).IP, conn.RemoteAddr().(*net.TCPAddr).Port)
		// go handleConnection(conn)
	}
}
