package controllers

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"github.com/yuvrajrathva/P2P-Gossips-Network/models"
)

func CreateSeed() {
	// Initialize Viper for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") 
	viper.AddConfigPath("../../config")   

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Parse the seed nodes from the configuration file
	seedNodes, err := parseSeedNodes()
	if err != nil {
		log.Fatalf("Error parsing seed nodes: %s", err)
	}

	callTCPServer(seedNodes)
	// Print the seed nodes
	// fmt.Println("Seed Nodes:")
	// printSeedNodes(seedNodes)
}

func parseSeedNodes() ([]models.ConfigSeed, error) {
	type Config struct {
		Seeds []models.ConfigSeed `mapstructure:"seeds"`
	}

	var config Config

	// Unmarshal the configuration file into the Config structure
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %s", err)
	}

	return config.Seeds, nil
}

// func printSeedNodes(seedNodes []models.ConfigSeed) {
// 	for _, seed := range seedNodes {
// 		fmt.Printf("IP: %s, Port: %d \n", seed.IP, seed.Port)
// 	}
// }

func getSeedNodes() ([]models.ConfigSeed, error) {
	// Initialize Viper for config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") 
	viper.AddConfigPath("../../config")   

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	type Config struct{
		Seeds []models.ConfigSeed
	}

	var config Config

	// Unmarshal the configuration file into the Config structure
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %s", err)
	}

	return config.Seeds, nil
}

func callTCPServer(seedNodes []models.ConfigSeed){
	var wg sync.WaitGroup
	wg.Add(len(seedNodes))
	fmt.Println("Starting Seed Nodes as Server")
	for _, seed := range seedNodes {
		go SeedTCPServer(seed.IP, seed.Port, &wg, &seed.PeerList)
	}
	wg.Wait()
}
