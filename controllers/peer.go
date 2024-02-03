package controllers

import (
	"fmt"
	"math/rand"

	"github.com/go-playground/locales/my"
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

	// Select a random n/2 + 1 seed nodes
	// n is the total number of seed nodes
	mySeedNodes := selectSeedNodes(seedNodeList)

	
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
