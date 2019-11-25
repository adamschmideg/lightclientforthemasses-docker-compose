package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

func check(maybeError error) {
	if maybeError != nil {
		fmt.Println(maybeError)
	}
}
func main() {
	server, _ := rpc.Dial("http://127.0.0.1:8545")
	client, _ := rpc.Dial("http://127.0.0.1:8546")
	type NodeInfo struct {
		ID string
	}
	var nodeInfo NodeInfo
	client.Call(&nodeInfo, "admin_nodeInfo")
	nodeID := nodeInfo.ID
	fmt.Println("nodeID", nodeID)

	var balance map[string]interface{}
	if err := server.Call(&balance, "les_addBalance", nodeID, 100000, "foobar"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(balance)
	}

	var clientInfo map[string]interface{}
	clientIDs := []string{nodeID}	
	if err := server.Call(&clientInfo, "les_clientInfo", clientIDs); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(clientInfo)
	}
}