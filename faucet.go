package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type balanceInfo struct {
	BalanceBefore int
	BalanceAfter int
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

type allClientsInfo map[string]map[string]interface{}
type clientInfo map[string]interface{}

func getBalance(serverRPCEndpoint string, nodeID string) (clientInfo,error) {
	var allInfo allClientsInfo
	var cInfo clientInfo
	log.Printf("GetBalance called at %s for %s", serverRPCEndpoint, nodeID)
	server, err := rpc.Dial(serverRPCEndpoint)
	if err != nil {
		return cInfo, fmt.Errorf("Server", err)
	}
	log.Printf("Server connected")

	clientIDs := []string{nodeID}	
	if err := server.Call(&allInfo, "les_clientInfo", clientIDs); err != nil {
		return cInfo, fmt.Errorf("clientinfo", err)
	}
	cInfo = allInfo[nodeID]
	log.Printf("Got balance %s", allInfo)
	return cInfo, nil
}

type formData struct {
	NodeID string
	Error error
	Balance balanceInfo
	Client clientInfo
	RPCEndpoint string
}

func makeRootHandler(rpcEndpoint string, templatePath string) func(http.ResponseWriter, *http.Request) {
	handler := func (w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		nodeID := r.FormValue("nodeID")
		var err error
		var bInfo balanceInfo
		var cInfo clientInfo
		log.Println("handle", r.Method, r.Form)

		switch {
		case nodeID == "":
			err = errors.New("nodeID is required")
		case r.Method == http.MethodPost:
			bInfo, err = addBalance(rpcEndpoint, nodeID, 1000, "foobar")
			cInfo, err = getBalance(rpcEndpoint, nodeID)
			if cInfo != nil {
				cInfo["pricing/oldBalance"] = bInfo.BalanceBefore
			}
		case r.Method == http.MethodGet:
			cInfo, err = getBalance(rpcEndpoint, nodeID)
		default:
			err = errors.New("Unsupported method")
		}

		fillData := formData{nodeID, err, bInfo, cInfo, rpcEndpoint}

		t, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Println("Parsing html", err)
			fmt.Fprintln(w, "Internal error")
			return
		}
		fmt.Println("fillData", fillData)
		if err := t.Execute(w, fillData); err != nil {
			log.Println("Executing template", err)
			fmt.Fprintln(w, "internal error")
		}
	}

	return handler
}

func lookupIP(address string) string {
	ips, _ := net.LookupIP(address)
	if (len(ips) > 0) {
		return ips[0].String()
	}
	return address
}

func main() {
	rpcaddr := flag.String("rpcaddr", "127.0.0.1", "Address of the lightserver's rpc endpoint")
	rpcport := flag.Int("rpcport", 8545, "Port of the lightserver's rpc endpoint")
	port := flag.Int("port", 8088, "Web service port of the faucet")
	templatePath := flag.String("template", "/var/www/faucet.html", "Full path to the html template file")
	flag.Parse()

	rpcIP := lookupIP(*rpcaddr)
	rpcEndpoint := fmt.Sprintf("http://%s:%v", rpcIP, *rpcport)
	rootHandler := makeRootHandler(rpcEndpoint, *templatePath)
	wsAddress := fmt.Sprintf(":%v", *port)
	log.Println("Listening at", wsAddress, ", calling", rpcEndpoint)
	
	http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(wsAddress, nil))
}