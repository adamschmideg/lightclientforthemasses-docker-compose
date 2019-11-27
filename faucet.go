package main

import (
	"fmt"
	//"html/template"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type balanceInfo struct {
	BalanceBefore int
	BalanceAfter int
	IsConnected bool
}

func addBalance(serverRPCEndpoint string, clientNodeID string, balance int, topic string) (balanceInfo,error) {
	var info balanceInfo
	server, err := rpc.Dial(serverRPCEndpoint)
	if err != nil {
		return info, fmt.Errorf("Server", err)
	}
	log.Printf("Server connected")

	var balances []int
	if err := server.Call(&balances, "les_addBalance", clientNodeID, balance, topic); err != nil {
		return info, fmt.Errorf("balance", err)
	} 
	info.BalanceBefore = balances[0]
	info.BalanceAfter = balances[1]
	log.Println("AddBalance success", balances)

	return info, nil
}

func getBalance(serverRPCEndpoint string, nodeID string) (balanceInfo,error) {
	var info balanceInfo
	log.Printf("GetBalance called at %s for %s", serverRPCEndpoint, nodeID)
	server, err := rpc.Dial(serverRPCEndpoint)
	if err != nil {
		return info, fmt.Errorf("Server", err)
	}
	log.Printf("Server connected")

	clientIDs := []string{nodeID}	
	if err := server.Call(&info, "les_clientInfo", clientIDs); err != nil {
		return info, fmt.Errorf("clientinfo", err)
	}
	log.Printf("Got balance %s", info)
	return info, nil
}

const rpcEndpoint string = "http://127.0.0.1:8545"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nodeID := r.FormValue("nodeID")
	if nodeID == "" {
		fmt.Fprint(w, "nodeID is required")
	} else if r.Method == "POST" {
		info, err := addBalance(rpcEndpoint, nodeID, 1000, "foobar")
		if err != nil {
			fmt.Fprintf(w, "Can't add balance %s", err)
		} else {
			fmt.Fprintf(w, "Added balance %s", info)
		}
	} else if r.Method == "GET" {
		info, err := getBalance(rpcEndpoint, nodeID)
		if err != nil {
			fmt.Fprintf(w, "Can't get info %s", err)
		} else {
			fmt.Fprintf(w, "info %s", info)
		}
	}
}

func _rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/*
	t, err := template.ParseFiles("test.html")
	if err != nil {
		log.Println("Problem", err)
	}
	fillData := struct{}{}
	*/

	nodeID := r.FormValue("nodeID")
	if nodeID == "" {
		fmt.Fprint(w, "nodeID is required")
	} else if r.Method == "POST" && r.FormValue("cmd") == "addBalance" {
	}
	if r.Method == "POST" {
		cmd := r.FormValue("cmd")
		if cmd == "addBalance" {

		}
		fmt.Fprintf(w, "addBalance(\"%s\"", r.FormValue("nodeID"))
	} else {
		//t.Execute(w, fillData)
		fmt.Fprintf(w, "<h1>Hello</h1>")
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(":8088", nil))
}