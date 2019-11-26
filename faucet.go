package main

import (
	"fmt"
	//"html/template"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type BalanceInfo struct {
	BalanceBefore int
	BalanceAfter int
	IsConnected bool
}

func AddBalance(serverRpcEndpoint string, clientNodeID string, balance int, topic string) (BalanceInfo,error) {
	var info BalanceInfo
	server, err := rpc.Dial(serverRpcEndpoint)
	if err != nil {
		return info, fmt.Errorf("Server", err)
	}

	var balances []int
	if err := server.Call(&balances, "les_addBalance", clientNodeID, balance, topic); err != nil {
		return info, fmt.Errorf("balance", err)
	} 
	info.BalanceBefore = balances[0]
	info.BalanceAfter = balances[1]

	return info, nil
}

func GetBalance(serverRpcEndpoint string, nodeID string) (BalanceInfo,error) {
	var info BalanceInfo
	server, err := rpc.Dial(serverRpcEndpoint)
	if err != nil {
		return info, fmt.Errorf("Server", err)
	}

	clientIDs := []string{nodeID}	
	if err := server.Call(&info, "les_clientInfo", clientIDs); err != nil {
		return info, fmt.Errorf("clientinfo", err)
	}
	return info, nil
}

const rpcEndpoint string = "http://127.0.0.1:8545"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nodeID := r.FormValue("nodeID")
	if nodeID == "" {
		fmt.Fprint(w, "nodeID is required")
	} else if r.Method == "POST" {
		info, err := AddBalance(rpcEndpoint, nodeID, 1000, "foobar")
		if err != nil {
			fmt.Fprintf(w, "Can't add balance %s", err)
		} else {
			fmt.Fprintf(w, "Added balance %s", info)
		}
	} else if r.Method == "GET" {
		info, err := GetBalance(rpcEndpoint, nodeID)
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