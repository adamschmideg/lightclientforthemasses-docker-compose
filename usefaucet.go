package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/rpc"
)

type nodeInfo struct {
	ID string
}

func useFaucet(clientRPCEndpoint string, faucetEndpoint string)(error) {

	rpcClient, err := rpc.Dial("http://127.0.0.1:8546")
	if err != nil {
		return fmt.Errorf("Client", err)
	}

	var nodeInfo nodeInfo
	if err := rpcClient.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		return fmt.Errorf("NodeInfo", err)
	}
	fmt.Println("Node", nodeInfo)

	formData := url.Values{
		"nodeID": {nodeInfo.ID},
	}
	resp, err := http.PostForm(faucetEndpoint, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	return nil
}

func main() {
	clientRPCEndpoint := "http://127.0.0.1:8546"
	faucetEndpoint := "http://127.0.0.1:8088"
	err := useFaucet(clientRPCEndpoint, faucetEndpoint)
	if err != nil {
		fmt.Println(err)
	}
}
