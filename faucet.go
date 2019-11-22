package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

const NodeID string = "fa25f829719067d8685332df3465bb877c2f49ee09587ae4afaa22623aaf3d7a"

func main() {
	client, _ := rpc.Dial("http://localhost:8545")

	type ServerInfo struct {
		FreeClientCapacity int
		MaximumCapacity int
		MinimumCapacity int
		PriorityConnectedCapacity int
		TotalCapacity int
		TotalConnectedCapacity int
	}

	if false {
		var dictResponse map[string]interface{}

		if err := client.Call(&dictResponse, "les_serverInfo"); err != nil {
			fmt.Println("Problem", err)
		}	else {
			fmt.Println("Response", dictResponse)
		}

		var arrayResponse []string

		if err := client.Call(&arrayResponse, "les_getCheckpoint", 1); err != nil {
			fmt.Println("Problem", err)
		}	else {
			fmt.Println("Response", arrayResponse)
		}
	}


	/* 
	// It hangs
	var dictResponse map[string]interface{}

	clientIDs := []string{NodeID}

	if err := client.Call(&dictResponse, "les_clientInfo", clientIDs); err != nil {
		fmt.Println("Problem", err)
	}	else {
		fmt.Println("Response", dictResponse)
	}

	var arrayResponse []string

	if err := client.Call(&arrayResponse, "les_addBalance", NodeID, 10000000, "qwerty"); err != nil {
		fmt.Println("Problem", err)
	}	else {
		fmt.Println("Response", arrayResponse)
	}
	*/
}