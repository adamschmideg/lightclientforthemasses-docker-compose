package main

import (
	"os/exec"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var rpcArgs = []string{"--rpc", "--rpcapi=admin,eth,les", "--rpcaddr=0.0.0.0"}

var port int = 30303
var rpcPort int = 8545

type geth struct {
	datadir string
	args []string
	cmd *exec.Cmd
}

func startGeth(datadir string, keepDatadir bool, args ...string) (*geth, error) {
	g := &geth{datadir, args, nil}
	if !keepDatadir {
		err := os.RemoveAll(datadir)
		if err != nil {
			return g, err
		}
	}
	allArgs := []string{"--datadir", datadir, "--port", fmt.Sprintf("%d", port), "--rpcport", fmt.Sprintf("%d", rpcPort)}
	port += 1
	rpcPort += 1
	allArgs = append(allArgs, rpcArgs...)
	allArgs = append(allArgs, args...)
	cmd := exec.Command("geth", allArgs...)
	log.Println("to start", cmd.String())
	err := cmd.Start()
	g.cmd = cmd
	if err != nil {
		return g, err
	}
	time.Sleep(1 * time.Second) // wait before we can attach to it
	return g, nil
}

func (g *geth) exec(js string) (string, error) {
	p, err := filepath.Abs(g.datadir)
	if err != nil {
		return "", err
	}
	ipcPath := filepath.Join(p, "geth.ipc")	
	execArgs := []string{"attach", "--exec", js, ipcPath}
	allArgs := append([]string{}, rpcArgs...)
	allArgs = append(allArgs, execArgs...)
	cmd := exec.Command("geth", allArgs...)
	log.Println("to exec", cmd.String())
	var b []byte
	b, err = cmd.CombinedOutput()
	out := strings.Trim(string(b), " \n\r\t\"")
	if err != nil {
		return out, err
	}
	return out, nil
}

func (g *geth) kill() error {
	err := g.cmd.Process.Kill()
	if err != nil {
		return err
	}
	_, err = g.cmd.Process.Wait()
	return err
}

func TestDemo(t *testing.T) {
	syncGoerli, err := startGeth("./datadirs/goerli/fast", true, "--goerli", "--light.serve=100", "--syncmode=fast", "--exitwhensynced")
	if err != nil {
		t.Fatal(syncGoerli.cmd.String(), err)
	}
	waited, err := syncGoerli.cmd.Process.Wait()
	if err != nil {
		t.Fatal("waiting for goerli to finish", err)
	}
	t.Log("fast sync finished", waited.String())
	server, err := startGeth("./datadirs/goerli/fast", true, "--light.serve=100", "--light.maxpeers=1", "--goerli", "--syncmode=fast")
	defer server.kill()
	if err != nil {
		t.Fatal(server.cmd.String(), err)
	}
	enode, err := server.exec("admin.nodeInfo.enode")
	if err != nil {
		t.Fatal(enode, err)
	}

	// Simple client
	client, err := startGeth("./datadirs/goerli/light", false, "--syncmode=light", "--nodiscover")
	defer client.kill()
	if err != nil {
		t.Fatal(client.cmd.String(), err)
	}
	clientRPC := rpc.Dial(ipcpath)
	client.Call("admin_addPeer")

	addPeerJs := fmt.Sprintf(`'admin.addPeer("%v"); admin.peers'`, enode)
	addPeerResult, err := client.exec(addPeerJs)
	if err != nil {
		t.Fatal(err, addPeerResult)
	}
	t.Log("addPeer", addPeerResult)
	time.Sleep(1 * time.Second) // wait before sync starts
	clientSyncResult, err := client.exec("eth.syncing")
	if clientSyncResult == "false" {
		t.Log("expected client to sync, but", clientSyncResult)
		t.Fail()
	}
	
	// Priority client
	prio, err := startGeth("./datadirs/goerli/prio", false, "--syncmode=light", "--nodiscover")
	defer prio.kill()
	if err != nil {
		t.Fatal(prio.cmd.String(), err)
	}
	addPeerToPrioResult, err := prio.exec(addPeerJs)
	if err != nil {
		t.Fatal(addPeerToPrioResult, err)
	}
	t.Log("addPeerToPrio", addPeerToPrioResult)
	nodeID, err := prio.exec("admin.nodeInfo.id")
	tokens := 3_000_000_000
	addBalanceJs := fmt.Sprintf(`'les.addBalance("%v", %v, "foobar")'`, nodeID, tokens)
	addBalanceResult, err := server.exec(addBalanceJs)
	if err != nil {
		t.Fatal("addBalance", addBalanceResult, err)
	}

	// Check if priority client is actually syncing and the regular client got kicked out
	time.Sleep(1 * time.Second) // wait before old client gets kicked out
	clientSyncResult, err = client.exec("eth.syncing")
	if clientSyncResult != "false" {
		t.Log("expected client sync to be false, but", clientSyncResult)
		t.Fail()
	}
	prioSyncResult, err := prio.exec("eth.syncing")	
	if prioSyncResult == "false" {
		t.Log("expected prio sync, but", prioSyncResult)
		t.Fail()
	}
}