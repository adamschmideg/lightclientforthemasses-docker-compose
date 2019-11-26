package main

import (
	"fmt"
	"net/http"

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

func addBalance()(error) {

	rpcClient, err := rpc.Dial("http://127.0.0.1:8546")
	if err != nil {
		return fmt.Errorf("Client", err)
	}

	var nodeInfo nodeInfo
	if err := rpcClient.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		return fmt.Errorf("NodeInfo", err)
	}
	fmt.Println("Node", nodeInfo)

	faucetEndpoint := "http://127.0.0.1:8088"
	url := fmt.Sprintf("%s?nodeID=%s", faucetEndpoint, nodeInfo.ID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func main() {
	err := addBalance()
	if err != nil {
		fmt.Println(err)
	}
}
