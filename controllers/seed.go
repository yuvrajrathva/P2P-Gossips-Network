package controllers

import (
	"fmt"
	"log"
	"sync"
	"os"
	"bufio"
	"strings"
	"strconv"

	"github.com/yuvrajrathva/P2P-Gossips-Network/models"
)

func CreateSeed() {
	// Parse the seed nodes from the configuration file
	outputFile, err := os.Create("../../outputfile.txt")
	if err != nil {
		log.Fatalf("Error creating output file: %s", err)
	}
	defer outputFile.Close()

	seedNodes, err := parseSeedNodes()
	if err != nil {
		log.Fatalf("Error parsing seed nodes: %s", err)
	}

	callTCPServer(seedNodes)  // start seed nodes as server
}

func parseSeedNodes() ([]models.ConfigSeed, error) {
	type Config struct {
		Seeds []models.ConfigSeed `mapstructure:"seeds"`
	}

	var config Config

	file, err := os.Open("../../config/config.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %s", err)
	}
	defer file.Close()

	// Create a new scanner to read the file
	scanner := bufio.NewScanner(file)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into IP and port
		ipPort := strings.Split(line, ":")
		ip := ipPort[0]
		port, err := strconv.Atoi(ipPort[1])

		if err != nil {
			return nil, fmt.Errorf("error converting port to int: %s", err)
		}

		// Add the seed node to the list
		config.Seeds = append(config.Seeds, models.ConfigSeed{IP: ip, Port: port})
	}

	return config.Seeds, nil
}

func getSeedNodes() ([]models.ConfigSeed, error) {
	// Parse the seed nodes from the configuration file
	seedNodes, err := parseSeedNodes()
	if err != nil {
		return nil, fmt.Errorf("error parsing seed nodes: %s", err)
	}

	return seedNodes, nil
}

func callTCPServer(seedNodes []models.ConfigSeed){
	var wg sync.WaitGroup
	wg.Add(len(seedNodes))
	fmt.Println("Starting Seed Nodes as Server")
	for _, seed := range seedNodes {
		go SeedTCPServer(seed.IP, seed.Port, &wg, &seed.PeerList) // start seed nodes as server
	}
	wg.Wait()
}
