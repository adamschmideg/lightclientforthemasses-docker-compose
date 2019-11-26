package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

func addBalance()(error) {
	type NodeInfo struct {
		ID string
	}

	server, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		return fmt.Errorf("Server", err)
	}
	client, err := rpc.Dial("http://127.0.0.1:8546")
	if err != nil {
		return fmt.Errorf("Client", err)
	}

	var nodeInfo NodeInfo
	if err := client.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		return fmt.Errorf("NodeInfo", err)
	}

	var balance []interface{}
	if err := server.Call(&balance, "les_addBalance", nodeInfo.ID, 100000, "foobar"); err != nil {
		return fmt.Errorf("balance", err)
	}
	fmt.Println(balance)

	var clientInfo []interface{}
	clientIDs := []string{nodeInfo.ID}	
	if err := server.Call(&clientInfo, "les_clientInfo", clientIDs); err != nil {
		return err
	}
	fmt.Println(clientInfo)
	return nil
}

func main() {
	err := addBalance()
	if err != nil {
		fmt.Println(err)
	}
}
