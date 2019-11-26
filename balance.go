package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

type nodeInfo struct {
	ID string
}

type balanceInfo struct {
	BalanceBefore int
	BalanceAfter int
	IsConnected bool
}

func addBalanceLocally(nodeID string, balance int, topic string) (balanceInfo, error) {
	var info balanceInfo
	server, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		return info, fmt.Errorf("Server", err)
	}

	var balances []int
	if err := server.Call(&balances, "les_addBalance", nodeID, balance, topic); err != nil {
		return info, fmt.Errorf("balance", err)
	} 
	info.BalanceBefore = balances[0]
	info.BalanceAfter = balances[1]

	//var clientInfo map[string]interface{}
	clientIDs := []string{nodeID}	
	if err := server.Call(&info, "les_clientInfo", clientIDs); err != nil {
		return info, fmt.Errorf("clientinfo", err)
	}
	return info, nil
}

func addBalance()(error) {

	client, err := rpc.Dial("http://127.0.0.1:8546")
	if err != nil {
		return fmt.Errorf("Client", err)
	}

	var nodeInfo nodeInfo
	if err := client.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		return fmt.Errorf("NodeInfo", err)
	}

	info, err := addBalanceLocally(nodeInfo.ID, 1000, "foobar")
	if err != nil {
		return err
	}
	fmt.Println(info)
	return nil
}

func main() {
	err := addBalance()
	if err != nil {
		fmt.Println(err)
	}
}
