package main

import (
	"os/exec"
	"testing"
)

func TestDemo(t *testing.T) {
	commonArgs := []string{"--rpc", "--rpcap=admin,eth,les", "--rpcaddr=0.0.0.0"}
// Start lightserver with 1 slot for light clients
	serverDatadir := "~/datadirs/goerli/fast"
	serverArgs :=  []string{"--lightserv=100", "--light.maxpeers=1", "--datadir", serverDatadir, "--goerli"}
	var args []string
	args = append([]string{}, commonArgs...)
	args = append(args, serverArgs...)
	cmdServer := exec.Command("geth", serverArgs...)
	if err := cmdServer.Start(); err != nil {
		t.Error("start", err)
	}
	if err := cmdServer.Process.Kill(); err != nil {
		t.Error("kill", err)
	}
// Get the server's enode
// Start a light client with an empty datadir
// Add the server as peer to let this client start syncing
// Start a priority client with an empty datadir
// Get the nodeID of the priority client
// Add balance for the priority client on the light server
// Add the server as peer to let priority client start syncing
// Check if it's actually syncing
// Check if the regular client got kicked out
}