package models

type ConfigSeed struct {
	IP   string
	Port int
	PeerList []Peer
}

type Peer struct {
	IP   string
	Port int
	MissedPings int 
	MessageList []MessageList
}

type MessageList struct{
	TimeStamp string
	MessageHash string
	IP string
	Port int
	PeerIP string
	PeerPort int
}
