package models

type PeerList struct {
	IP   string
	Port int
}

type ConfigSeed struct {
	IP   string
	Port int
	PeerList []PeerList
}

type Peer struct {
	IP   string
	Port int
}
