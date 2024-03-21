package utils

type Node struct {
	Id      int
	Address string
	Leader  int // contiene l'ID del leader
}

type PeerAddr struct {
	PeerAddress string
}

type NodeList struct {
	Nodes []Node
}

type NodeINFO struct {
	Id      int
	Address string
	List    NodeList
	Leader  int
}

type Message struct {
	SkipCount int
	MexID     int
	CurrNode  NodeINFO
}
