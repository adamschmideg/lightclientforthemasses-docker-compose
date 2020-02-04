package main

import (
	"os/exec"
	"testing"
	"time"
)

func TestEnode(t *testing.T) {
	var err error
	serverDatadir := "~/datadirs/goerli/fast"
	// Get the server's enode
	enodeExec := exec.Command("geth", "--datadir", serverDatadir, "attach", "--exec", "admin.nodeInfo.enode" )
	var out []byte
	out, err = enodeExec.CombinedOutput()
	if err != nil {
		t.Error("enode", err)
	}
	enode := string(out)
	t.Log("yes", enode)
}

func TestDemo(t *testing.T) {
	var err error
	commonArgs := []string{"--rpc", "--rpcap=admin,eth,les", "--rpcaddr=0.0.0.0"}
	// Start lightserver with 1 slot for light clients
	serverDatadir := "./datadirs/goerli/fast"
	serverArgs :=  []string{"--lightserv=100", "--light.maxpeers=1", "--datadir", serverDatadir, "--goerli"}
	var args []string
	args = append([]string{}, commonArgs...)
	args = append(args, serverArgs...)
	cmdServer := exec.Command("geth", serverArgs...)
	if err = cmdServer.Start(); err != nil {
		t.Error("start", err)
	}

	// Get the server's enode
	time.Sleep(1 * time.Second) // wait before we can attach to it
	enodeExec := exec.Command("geth", "attach", "--exec", "admin.nodeInfo.enode", "./datadirs/goerli/fast/geth.ipc")
	var out []byte
	out, err = enodeExec.CombinedOutput()
	if err != nil {
		t.Error("enode", err, string(out))
		out = nil
	}
	enode := string(out)
	t.Log("yes", enode)

	// Tear down: kill all process
	if err = cmdServer.Process.Kill(); err != nil {
		t.Error("kill", err)
	}

// Start a light client with an empty datadir
// Add the server as peer to let this client start syncing
// Start a priority client with an empty datadir
// Get the nodeID of the priority client
// Add balance for the priority client on the light server
// Add the server as peer to let priority client start syncing
// Check if it's actually syncing
// Check if the regular client got kicked out
}