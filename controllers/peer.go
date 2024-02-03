package controllers

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/yuvrajrathva/P2P-Gossips-Network/models"
)

// CreatePeer initializes a peer with the specified IP and port
func CreatePeer(ip string, port int) {
	peer := &models.Peer{
		IP:   ip,
		Port: port,
	}
	
	seedNodeList, err := getSeedNodes()
	if err != nil {
		fmt.Printf("Error getting seed nodes: %s\n", err)
		return
	}

	selectedSeedNodes := selectSeedNodes(seedNodeList) // select n/2 + 1 seed nodes

	// var selectedPeers []models.Peer
	selectedPeersList := selectPeersList(selectedSeedNodes) // select all peers from selected seed nodes
	selectedPeers := selectPeers(selectedPeersList) // select 4 peers from selected peers list
	addPeerToSeedNodes(selectedSeedNodes, peer) // add this peer to Peer List of selected seed nodes

	if(len(selectedPeers) == 0){
		TCPServer(ip, port) // start server
	}else{
		startServer(selectedPeers) // start all selected peers as a server
		connectToPeersServer(selectedPeers) // connect to selected peers
	}
	fmt.Printf("Peer created - IP: %s, Port: %d\n", peer.IP, peer.Port)
}

func selectSeedNodes(seedNodes []models.ConfigSeed) []models.ConfigSeed {
	n := len(seedNodes)
	numSeedNodes := n/2 + 1
	selectedSeedNodes := make([]models.ConfigSeed, numSeedNodes)

	// Shuffle the seed nodes
	rand.Shuffle(n, func(i, j int) {
		seedNodes[i], seedNodes[j] = seedNodes[j], seedNodes[i]
	})

	// Select n/2 + 1 seed nodes
	for i := 0; i < numSeedNodes; i++ {
		selectedSeedNodes[i] = seedNodes[i]
	}

	return selectedSeedNodes
}

func selectPeersList(seedNodes []models.ConfigSeed) []models.Peer {
	var selectedPeers []models.Peer
	// we do not append duplicate peer in selectedPeers

	for _, seed := range seedNodes {
		for _, p := range seed.PeerList {
			if !containsPeer(selectedPeers, p) {
				selectedPeers = append(selectedPeers, p)
			}
		}
	}

	return selectedPeers
}

func containsPeer(peers []models.Peer, peer models.Peer) bool {
	for _, p := range peers {
		if p == peer {
			return true
		}
	}

	return false
}

func selectPeers(peers []models.Peer) []models.Peer {
	n := len(peers)
	if(n < 4){
		return peers
	}
	selectedPeers := make([]models.Peer, 4)

	// Shuffle the peers
	rand.Shuffle(n, func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})

	for i := 0; i < 4; i++ {
		selectedPeers[i] = peers[i]
	}

	return selectedPeers
}

func addPeerToSeedNodes(seedNodes []models.ConfigSeed, peer *models.Peer) {
	for i := range seedNodes {
		seedNodes[i].PeerList = append(seedNodes[i].PeerList, *peer)
	}
}

func startServer(peers []models.Peer) {
	for _, p := range peers {
		go TCPServer(p.IP, p.Port)
	}
}

func connectToPeersServer(peers []models.Peer) {
	for _, p := range peers {
		fmt.Printf("Connecting to peer - IP: %s, Port: %d\n", p.IP, p.Port)
		go TCPClient(p.IP, p.Port) // connect to peers server
	}
}

func TCPClient(ip string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error connecting to peer - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to peer - IP: %s, Port: %d\n", ip, port)
}

func TCPServer(ip string, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Error starting server - IP: %s, Port: %d, Error: %s\n", ip, port, err)
		return
	}
	defer listener.Close()

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
