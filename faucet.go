package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

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
	var resp ServerInfo
	//var resp map[string]interface{}

	if err := client.Call(&resp, "les_serverInfo"); err != nil {
		fmt.Println("Problem", err)
	}	else {
		fmt.Println("Response", resp)
	}
}