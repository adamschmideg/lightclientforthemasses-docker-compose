package main

import (
	"os/exec"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestEnode(t *testing.T) {
	enode := "  foobar   "
	addPeer := fmt.Sprintf(`'admin.addPeer("%v")'`, strings.TrimSpace(enode))
	t.Log(addPeer)
}

func TestDemo(t *testing.T) {
	var err error
	commonArgs := []string{"--rpc", "--rpcapi=admin,eth,les", "--rpcaddr=0.0.0.0"}
	// Start lightserver with 1 slot for light clients
	serverDatadir := "./datadirs/goerli/fast"
	serverArgs :=  []string{"--lightserv=100", "--light.maxpeers=1", "--datadir", serverDatadir, "--goerli"}
	var args []string
	args = append([]string{}, commonArgs...)
	args = append(args, serverArgs...)
	cmdServer := exec.Command("geth", args...)
	if err = cmdServer.Start(); err != nil {
		t.Error("start", err)
	}
	t.Log("started", cmdServer.String())
	time.Sleep(1 * time.Second) // wait before we can attach to it

	// Get the server's enode
	enodeExec := exec.Command("geth", "attach", "--exec", "admin.nodeInfo.enode", "--datadir", serverDatadir)
	var out []byte
	out, err = enodeExec.CombinedOutput()
	if err != nil {
		t.Error("enode", err, string(out))
		out = nil
	}
	enodeRaw := string(out)
	enode := strings.Trim(enodeRaw, " \n\r\"")

	// Start a light client with an empty datadir
	clientDir := "./datadirs/goerli/light"
	err = os.RemoveAll(clientDir)
	if err != nil {
		t.Error("rmdir", err)
	}
	clientArgs := []string{"--syncmode=light", "--datadir", clientDir, "--nodiscover"}
	args = append([]string{}, commonArgs...)
	args = append(args, clientArgs...)
	cmdClient := exec.Command("geth", args...)
	if err = cmdClient.Start(); err != nil {
		t.Error("start client", err)
	}
	t.Log("started", cmdClient.String())
	time.Sleep(1 * time.Second) // wait before we can attach to it

	// Add the server as peer to let this client start syncing
	addPeerJs := fmt.Sprintf(`'admin.addPeer("%v")'`, enode)
	t.Log("peer", addPeerJs)
	cmdAddPeer := exec.Command("geth", "attach", "--datadir", clientDir, "--exec", addPeerJs)
	out, err = cmdAddPeer.CombinedOutput()
	if err != nil {
		t.Error("addPeer", cmdAddPeer.String(), err, string(out))
		out = nil
	}
	addPeerResult := string(out)
	t.Log("addPeer", addPeerResult)

	// Tear down: kill all process
	if err = cmdServer.Process.Kill(); err != nil {
		t.Error("kill server", err)
	}
	if err = cmdClient.Process.Kill(); err != nil {
		t.Error("kill client", err)
	}

// Start a priority client with an empty datadir
// Get the nodeID of the priority client
// Add balance for the priority client on the light server
// Add the server as peer to let priority client start syncing
// Check if it's actually syncing
// Check if the regular client got kicked out
}