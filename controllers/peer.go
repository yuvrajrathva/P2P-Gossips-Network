package controllers

import (
	"fmt"
	"math/rand"
	"sync"
	"os"

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
	fmt.Println("Number of total Seed Nodes:", len(seedNodeList))

	selectedSeedNodes := selectSeedNodes(seedNodeList) // select n/2 + 1 seed nodes

	selectedPeersList := selectPeersList(selectedSeedNodes) // select all peers from selected seed nodes
	selectedPeers := selectPeers(selectedPeersList)         // select 4 peers from selected peers list
	printPeerNodes(selectedPeers)

	if len(selectedPeers) > 0 {
		connectToPeersServer(selectedPeers) // connect to selected peers
	}

	addPeerToSeedNodes(selectedSeedNodes, peer) // add this peer to Peer List of selected seed nodes
	var wg sync.WaitGroup
	wg.Add(1)
	go PeerTCPServer(ip, port, &wg, &selectedPeers) // start peer server
	wg.Wait()
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
	for i := range seedNodes {
		requestingPeerList(seedNodes[i].IP, seedNodes[i].Port, &seedNodes[i].PeerList)
	}
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
		if p.IP == peer.IP && p.Port == peer.Port {
			return true
		}
	}
	return false
}

func selectPeers(peers []models.Peer) []models.Peer {
	n := len(peers)
	if n < 4 {
		return peers
	}
	selectedPeers := make([]models.Peer, 4)

	// Shuffle the peers
	rand.Shuffle(n, func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})

	// select 4 peers
	for i := 0; i < 4; i++ {
		selectedPeers[i] = peers[i]
	}

	return selectedPeers
}

func addPeerToSeedNodes(seedNodes []models.ConfigSeed, peer *models.Peer) {
	for _, seed := range seedNodes {
		PeerTCPClient(seed.IP, seed.Port, peer) // add peer to seed node
	}
}

func connectToPeersServer(peers []models.Peer) {
	for _, p := range peers {
		fmt.Printf("Connecting... to peer - %s:%d\n", p.IP, p.Port)
		go TCPClient(p.IP, p.Port) // connect to peer
	}
}

func printPeerNodes(peerNodes []models.Peer) {
	outputfile, err := os.OpenFile("../../outputfile.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file: ", err.Error())
		return
	}
	defer outputfile.Close()

	for _, peer := range peerNodes {
		fmt.Printf("Selected Peer IP: %s, Port: %d\n", peer.IP, peer.Port)
		outputfile.WriteString(fmt.Sprintf("Selected Peer IP: %s, Port: %d\n", peer.IP, peer.Port))
	}
}
