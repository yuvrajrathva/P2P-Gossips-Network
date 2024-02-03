package models

type ConfigSeed struct {
	IP   string
	Port int
	PeerList []Peer
}

type Peer struct {
	IP   string
	Port int
}
