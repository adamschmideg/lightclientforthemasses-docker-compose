
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
    "testing"
)

func TestThrottling(t *testing.T) {
	wsaddr := flag.String("wsaddr", "127.0.0.1", "Address of the faucet")
	flag.Parse()
	ip, _ := net.LookupIP(*wsaddr)
	endpoint := fmt.Sprintf("http://%s:%v", ip, "8088")
	
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