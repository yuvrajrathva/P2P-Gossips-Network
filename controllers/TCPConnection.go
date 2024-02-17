package controllers

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yuvrajrathva/P2P-Gossips-Network/models"
)

func TCPClient(ip string, port int) { // connect to peer server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error connecting to peer - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to peer - IP: %s, Port: %d\n", ip, port)
}

func PeerTCPClient(ip string, port int, peer *models.Peer) { // connect to seed server and send this peer details to register
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

func SeedTCPServer(ip string, port int, wg *sync.WaitGroup, peerList *[]models.Peer) { // start seed server
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

	// this go routine will check liveness of all peers in peerList
	go func() {
		for {
			var wg sync.WaitGroup
			for _, peer := range *peerList {
				wg.Add(1)
				go PeerLivelinessChecker(ip, port, peer.IP, peer.Port, &wg, &peer.MissedPings, peerList)
			}
			wg.Wait()
			time.Sleep(13 * time.Second) // check liveness of all peers in peerList after every 13 seconds
		}
	}()

	// this go routine will send gossip message to all peers in peerList
	go func() {
		var wg sync.WaitGroup
		for _, peer := range *peerList {
			wg.Add(1)
			go GossipMessage(peer.IP, peer.Port, ip, port, &wg)
		}
		wg.Wait()
		time.Sleep(5 * time.Second) // send gossip message to all peers in peerList after every 5 seconds
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

func requestingPeerList(ip string, port int, peerList *[]models.Peer) { // request for peer list from seed server
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

	*peerList = stringToPeerList(string(buffer[:n])) // update peer list
}

func handleSeedServerConnection(conn net.Conn, peerList *[]models.Peer) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return
	}

	outputfile, err := os.OpenFile("../../outputfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file: ", err.Error())
		return
	}
	defer outputfile.Close()

	if string(buffer[:n]) == "peerList\n" { // send peer list to peer
		_, err = conn.Write([]byte(peerListToString(peerList)))
		if err != nil {
			fmt.Println("Error while sending peer list: ", err.Error())
			return
		}
	} else if stringToArray(string(buffer[:n]), ":")[0] == "removePeer" { // remove peer from peer list
		deadIP := stringToArray(string(buffer[:n]), ":")[1]
		deadPort, _ := strconv.Atoi(stringToArray(string(buffer[:n]), ":")[2])

		for i, peer := range *peerList {
			if peer.IP == deadIP && peer.Port == deadPort {
				*peerList = append((*peerList)[:i], (*peerList)[i+1:]...)

				fmt.Printf("Dead Node: %s:%d:%s:%s\n", deadIP, deadPort, time.Now().String(), conn.LocalAddr().(*net.TCPAddr).IP.String()+":"+strconv.Itoa(conn.LocalAddr().(*net.TCPAddr).Port))

				outputfile.WriteString(fmt.Sprintf("Dead Node: %s:%d:%s:%s\n", deadIP, deadPort, time.Now().String(), conn.LocalAddr().(*net.TCPAddr).IP.String()+":"+strconv.Itoa(conn.LocalAddr().(*net.TCPAddr).Port)))
				break
			}
		}
	} else { // add peer to peer list
		peer := stringToArray(string(buffer[:n]), ":")
		ip := peer[0]
		port, err := strconv.Atoi(peer[1])
		if err != nil {
			fmt.Println("Error while converting port to int: ", err.Error())
			return
		}
		*peerList = append(*peerList, models.Peer{IP: ip, Port: port})

		fmt.Printf("Peer added - IP: %s, Port: %d to Seed: %s \n", ip, port, conn.LocalAddr().(*net.TCPAddr))

		outputfile.WriteString(fmt.Sprintf("Peer added - IP: %s, Port: %d to Seed: %s \n", ip, port, conn.LocalAddr().(*net.TCPAddr).IP.String()+":"+strconv.Itoa(conn.LocalAddr().(*net.TCPAddr).Port)))
	}
}

func peerListToString(peerList *[]models.Peer) string { // convert peer list to string
	var str string
	for _, peer := range *peerList {
		str += fmt.Sprintf("IP: %s, Port: %d\n", peer.IP, peer.Port)
	}
	return str
}

func stringToPeerList(str string) []models.Peer {  // convert string to peer list
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

func stringToArray(str string, sep string) []string { // convert string to array using separator
	return strings.Split(str, sep)
}

func handlePeerServerConnection(conn net.Conn, ip string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error while reading: ", err.Error())
		return
	}

	if stringToArray(string(buffer[:n]), " : ")[0] == "GossipMessage" { // handle gossip message
		messageTimeStamp := stringToArray(string(buffer[:n]), " : ")[1]
		senderIP := stringToArray(string(buffer[:n]), " : ")[2]
		senderPort:=stringToArray(string(buffer[:n]), " : ")[3]
		messageHash := stringToArray(string(buffer[:n]), " : ")[4]

		outputfile, err := os.OpenFile("../../outputfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening output file: ", err.Error())
			return
		}
		defer outputfile.Close()
		fmt.Printf("Gossip Message Received: %s : %s : %s\n", messageTimeStamp, senderIP+":"+senderPort, messageHash)
		outputfile.WriteString(fmt.Sprintf("Gossip Message Received: %s : %s : %s\n", messageTimeStamp, senderIP+":"+senderPort, messageHash))
	} else { // handle liveness message
		senderTimestamp := stringToArray(string(buffer[:n]), ":")[1]
		senderIP := stringToArray(string(buffer[:n]), ":")[2]

		_, err = conn.Write([]byte(senderTimestamp + ":" + senderIP + ":" + ip + ":\n"))
		if err != nil {
			fmt.Println("Error while sending liveness message: ", err.Error())
			return
		}
	}
}

func PeerLivelinessChecker(selfIP string, selfPort int, ip string, port int, wg *sync.WaitGroup, missedPings *int, peerList *[]models.Peer) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		*missedPings = *missedPings + 1
		if *missedPings >= 3 { // remove peer from seed if missedPings >= 3
			fmt.Printf("Peer is dead - IP: %s, Port: %d\n", ip, port)
			removePeerFromSeedNodes(selfIP, selfPort, ip, port, peerList)
		}

		// fmt.Printf("Error connecting to peer - IP: %s, Port: %d, Missed Pings: %d, Error: %s\n", ip, port, *missedPings, err)
		removePeerFromSeedNodes(selfIP, selfPort, ip, port, peerList)
	} else { // send liveness message to peers
		defer conn.Close()

		_, err = conn.Write([]byte("LivenessMessage:" + time.Now().String() + ":" + selfIP))
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
			removePeerFromSeedNodes(selfIP, selfPort, ip, port, peerList)
			return
		}
		*missedPings = 0
		fmt.Printf("Peer is alive - IP: %s, Port: %d\n", ip, port)
	}
	defer wg.Done()
}

func removePeerFromSeedNodes(selfIP string, selfPort int, ip string, port int, peerList *[]models.Peer) { // remove dead peer from seed nodes
	seedNodeList, err := getSeedNodes()
	if err != nil {
		fmt.Printf("Error getting seed nodes: %s\n", err)
		return
	}

	for _, seed := range seedNodeList {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", seed.IP, seed.Port))
		if err != nil {
			fmt.Printf("Error connecting to seed server - IP: %s, Port: %d, Error: %s\n", seed.IP, seed.Port, err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte("removePeer:" + ip + ":" + strconv.Itoa(port)))
		if err != nil {
			fmt.Println("Error while sending remove peer message:", err)
			return
		}
	}

	outputfile, err := os.OpenFile("../../outputfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file: ", err.Error())
		return
	}
	defer outputfile.Close()

	// remove dead peer from peer list
	for i, peer := range *peerList {
		if peer.IP == ip && peer.Port == port {
			*peerList = append((*peerList)[:i], (*peerList)[i+1:]...)
			fmt.Printf("Dead Node: %s:%d:%s:%s\n", ip, port, time.Now().String(), selfIP+":"+strconv.Itoa(selfPort))

			outputfile.WriteString(fmt.Sprintf("Dead Node: %s:%d:%s:%s\n", ip, port, time.Now().String(), selfIP+":"+strconv.Itoa(selfPort)))
			break
		}
	}
}

func GossipMessage(peerIP string, peerPort int, selfIP string, selfPort int, wg *sync.WaitGroup) { // send gossip message to peers
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", peerIP, peerPort))
	if err != nil {
		fmt.Printf("Error connecting to peer for gossip message- IP: %s, Port: %d, Error: %s\n", peerIP, peerPort, err)
		return
	}
	defer conn.Close()

	Message := "Gossip Message"
	hash := sha256.New()
	hash.Write([]byte(Message))
	hashedMessage := hash.Sum(nil)

	// send gossip message to peers
	_, err = conn.Write([]byte("GossipMessage : " + time.Now().String() + " : " + selfIP + " : " + strconv.Itoa(selfPort) + " : " + string(hashedMessage)))
	if err != nil {
		fmt.Println("Error while sending gossip message:", err)
		return
	}

	defer wg.Done()
}
