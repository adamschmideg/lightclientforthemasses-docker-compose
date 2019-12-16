
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
    "testing"
)


var wsaddr = flag.String("wsaddr", "127.0.0.1", "Address of the faucet")
func Init() {
	flag.Parse()
}

func TestThrottling(t *testing.T) {
	ips, err := net.LookupIP(*wsaddr)
	if len(ips) == 0 || err != nil {
		t.Error("Cannot resolve", *wsaddr)
		return
	}
	ip := ips[0].String()
	endpoint := fmt.Sprintf("http://%s:%v", ip, "8088")
	fmt.Println("debug", ip, endpoint)
	
	resp1, err := http.Get(endpoint)
	if err != nil {
		t.Errorf("first %s", err)
	}
	defer resp1.Body.Close()
	resp2, err := http.Get(endpoint)
	if resp2.StatusCode == http.StatusOK {
		t.Errorf("Expected error for second %s", resp2.Status)
	}
	defer resp2.Body.Close()
}