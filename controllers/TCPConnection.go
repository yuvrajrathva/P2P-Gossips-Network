package controllers

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yuvrajrathva/P2P-Gossips-Network/models"
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

func PeerTCPClient(ip string, port int, peer *models.Peer) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error connecting to seed - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(peer.IP + ":" + strconv.Itoa(peer.Port)))
	if err != nil {
		fmt.Println("Error while sending peer details:", err)
		return
	}
}

func SeedTCPServer(ip string, port int, wg *sync.WaitGroup, peerList *[]models.Peer) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error starting seed server - IP: %s, Port: %d, Error: %s\n", ip, port, err)
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

		go handleSeedServerConnection(conn, peerList)
	}
}

func PeerTCPServer(ip string, port int, wg *sync.WaitGroup, peerList *[]models.Peer) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error starting peer server - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer listener.Close()
	defer wg.Done()

	fmt.Printf("Server started - IP: %s, Port: %d\n", ip, port)

	go func() {
		for {
			var wg sync.WaitGroup
			for _, peer := range *peerList {
				wg.Add(1)
				go PeerLivelinessChecker(ip, port, peer.IP, peer.Port, &wg, &peer.MissedPings)
			}
			wg.Wait()
			time.Sleep(13 * time.Second)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err)
			return
		}

		defer conn.Close()

		// fmt.Println("Peer connected:", port)
		go handlePeerServerConnection(conn, ip)
	}
}

func requestingPeerList(ip string, port int, peerList *[]models.Peer) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error connecting to seed server - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("peerList\n"))
	if err != nil {
		fmt.Println("Error while requesting for peer list:", err)
		return
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading peer list:", err)
		return
	}

	*peerList = stringToPeerList(string(buffer[:n]))

	// fmt.Printf("Peer List from Seed server %s:%d is %v\n", ip, port, *peerList)
}

func handleSeedServerConnection(conn net.Conn, peerList *[]models.Peer) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return
	}

	if string(buffer[:n]) == "peerList\n" {
		// send peer list
		// fmt.Printf("I m inside peerList request");
		// fmt.Printf("Peer List: %v\n", peerList)
		_, err = conn.Write([]byte(peerListToString(peerList)))
		if err != nil {
			fmt.Println("Error while sending peer list: ", err.Error())
			return
		}
	} else {
		peer := stringToArray(string(buffer[:n]), ":")
		ip := peer[0]
		port, err := strconv.Atoi(peer[1])
		if err != nil {
			fmt.Println("Error while converting port to int: ", err.Error())
			return
		}
		*peerList = append(*peerList, models.Peer{IP: ip, Port: port})
		fmt.Printf("Peer added - IP: %s, Port: %d to Seed: %s \n", ip, port, conn.LocalAddr().(*net.TCPAddr))

		// getPeerListFromSeeds()
	}
}

func peerListToString(peerList *[]models.Peer) string {
	var str string
	for _, peer := range *peerList {
		str += fmt.Sprintf("IP: %s, Port: %d\n", peer.IP, peer.Port)
	}
	return str
}

func stringToPeerList(str string) []models.Peer {
	var peerList []models.Peer
	peerListStr := string(str)
	peerListStr = peerListStr[:len(peerListStr)-1]
	peerListStrArr := stringToArray(peerListStr, "\n")
	for _, peerStr := range peerListStrArr {
		peer := stringToArray(peerStr, ", ")
		ip := stringToArray(peer[0], ": ")[1]
		port := stringToArray(peer[1], ": ")[1]
		portInt, _ := strconv.Atoi(port)
		peerList = append(peerList, models.Peer{IP: ip, Port: portInt})
	}
	return peerList
}

func stringToArray(str string, sep string) []string {
	return strings.Split(str, sep)
}

// func getPeerListFromSeeds() {
// 	seedNodeList, err := getSeedNodes()
// 	if err != nil {
// 		fmt.Printf("Error getting seed nodes: %s\n", err)
// 		return
// 	}

// 	for _, seed := range seedNodeList {
// 		requestingPeerList(seed.IP, seed.Port, &seed.PeerList)
// 	}
// }

func handlePeerServerConnection(conn net.Conn, ip string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return
	}

	senderTimestamp := stringToArray(string(buffer[:n]), ":")[0]
	senderIP := stringToArray(string(buffer[:n]), ":")[1]

	_, err = conn.Write([]byte(senderTimestamp + ":" + senderIP + ":" + ip + ":\n"))
	if err != nil {
		fmt.Println("Error while sending liveness message: ", err.Error())
		return
	}
}

func PeerLivelinessChecker(selfIP string, selfPort int, ip string, port int, wg *sync.WaitGroup, missedPings *int) {
	// detect dead peers and remove that peer node from seed if missedPings >= 3
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		*missedPings = *missedPings + 1
		if *missedPings >= 3 {
			fmt.Printf("Peer is dead - IP: %s, Port: %d\n", ip, port)
			removePeerFromSeedNodes(ip, port)
		}

		fmt.Printf("Error connecting to peer - IP: %s, Port: %d, Missed Pings: %d, Error: %s\n", ip, port, *missedPings, err)
	} else {
		defer conn.Close()

		_, err = conn.Write([]byte(time.Now().String() + ":" + selfIP + "\n"))
		if err != nil {
			fmt.Println("Error while sending liveness message: ", err.Error())
			return
		}

		buffer := make([]byte, 1024)

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error while reading liveness message:", err)
			return
		}

		if n == 0 {
			fmt.Printf("Peer is dead - IP: %s, Port: %d\n", ip, port)
			removePeerFromSeedNodes(ip, port)
			return
		}
		*missedPings = 0
		fmt.Printf("Peer is alive - IP: %s, Port: %d\n", ip, port)
	}
	defer wg.Done()
}

func removePeerFromSeedNodes(ip string, port int) {
	seedNodeList, err := getSeedNodes()
	if err != nil {
		fmt.Printf("Error getting seed nodes: %s\n", err)
		return
	}

	for _, seed := range seedNodeList {
		for i, peer := range seed.PeerList {
			if peer.IP == ip && peer.Port == port {
				seed.PeerList = append(seed.PeerList[:i], seed.PeerList[i+1:]...)
				break
			}
		}
	}

}
