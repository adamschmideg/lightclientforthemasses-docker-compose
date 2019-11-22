package main

import (
	"fmt"
	//"html/template"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

const NodeID string = "fa25f829719067d8685332df3465bb877c2f49ee09587ae4afaa22623aaf3d7a"

func connectToLightServer() {
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	/*
	t, err := template.ParseFiles("test.html")
	if err != nil {
		log.Println("Problem", err)
	}
	fillData := struct{}{}
	*/

	if r.Method == "POST" {
		fmt.Fprintf(w, "addBalance(\"%s\"", r.FormValue("nodeID"))
	} else {
		//t.Execute(w, fillData)
		fmt.Fprintf(w, "<h1>Hello</h1>")
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}