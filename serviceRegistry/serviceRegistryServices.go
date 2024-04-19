package main

import (
	"SDCCproject/utils"

	"log"
	"net/rpc"
)

// NodeHandler is the interface which exposes the RPC method ManageNode
type NodeHandler struct{}

// ManageNode is the discovery service for each new DockerfilePeer in DS
func (NodeHandler) ManageNode(peerAddress utils.PeerAddr, nodeInfo *utils.NodeINFO) error {
	nodes := nodeList.GetAllNodes()

	for _, node := range nodes {
		/* Check if the node was already in the system */
		if peerAddress.PeerAddress == node.Address {
			nodeInfo.Id = node.Id
			nodeInfo.List.Nodes = nodeList.GetAllNodes()
			nodeInfo.Leader = node.Leader
			nodeInfo.Participant = node.Participant

			return nil
		}
	}

redo:
	newID := utils.Random(0, 20)
	for _, node := range nodes {
		/* Check if ID was already assigned */
		if newID == node.Id {
			goto redo
		}
	}

	var newNode utils.Node
	newNode = utils.Node{Id: newID, Address: peerAddress.PeerAddress, Leader: -1, Participant: false}

	nodeList.AddNode(newNode)

	/* Send Peer info */
	nodeInfo.Id = newID
	nodeInfo.List.Nodes = nodeList.GetAllNodes()
	nodeInfo.Leader = newNode.Leader
	nodeInfo.Participant = newNode.Participant // setting false because new node in system is initially set as non-participant

	/* Update list to the other node in DS */
	for _, node := range nodes {
		peer, err := rpc.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		err = peer.Call("PeerServiceHandler.UpdateList", newNode, nil)
		if err != nil {
			log.Fatal("List update failed: ", err)
		}

		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}

	return nil
}
