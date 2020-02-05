package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

var rpcArgs = []string{"--rpc", "--rpcapi=admin,eth,les"}

var port int = 30303
var rpcPort int = 8545

type geth struct {
	datadir string
	args    []string
	cmd     *exec.Cmd
	rpc     *rpc.Client
}

func startGeth(datadir string, keepDatadir bool, args ...string) (*geth, error) {
	g := &geth{datadir, args, nil, nil}
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
	var err error
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	g.cmd = cmd
	time.Sleep(1 * time.Second) // wait before we can attach to it
	// TODO: probe for it properly
	g.rpc, err = rpc.Dial(g.ipcpath())
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *geth) ipcpath() string {
	return filepath.Join(g.datadir, "geth.ipc")
}

func (g *geth) kill() error {
	err := g.cmd.Process.Kill()
	if err != nil {
		return err
	}
	_, err = g.cmd.Process.Wait()
	return err
}

func (g *geth) addPeer(enode string) error {
	peerCh := make(chan *p2p.PeerEvent)
	sub, err := g.rpc.Subscribe(context.Background(), "admin", peerCh, "peerEvents")
	if err != nil {
		return fmt.Errorf("subscribe: %v", err)
	}
	defer sub.Unsubscribe()
	if err := g.rpc.Call(nil, "admin_addPeer", enode); err != nil {
		return fmt.Errorf("admin_addPeer: %v", err)
	}
	select {
	case ev := <-peerCh:
		fmt.Print("event", ev)

	case err := <-sub.Err():
		return fmt.Errorf("notification: %v", err)
	}
	return nil
}

func (g *geth) waitSynced() error {
	ch := make(chan interface{})
	sub, err := g.rpc.Subscribe(context.Background(), "eth", ch, "syncing")
	if err != nil {
		return fmt.Errorf("syncing: %v", err)
	}
	defer sub.Unsubscribe()
	timeout := time.After(40 * time.Second)
	for {
		select {
		case ev := <-ch:
			syncing, ok := ev.(bool)
			if ok && !syncing {
				return nil
			}
		case err := <-sub.Err():
			return fmt.Errorf("notification: %v", err)
		case <-timeout:
			return fmt.Errorf("timeout syncing")
		}
	}
}

func TestDemo(t *testing.T) {
	server, err := startGeth("./datadirs/goerli/fast", true, "--light.serve=100", "--light.maxpeers=1", "--goerli", "--syncmode=fast", "--nat=extip:127.0.0.1")
	defer server.kill()
	if err != nil {
		t.Fatal(server.cmd.String(), err)
	}
	if err := server.waitSynced(); err != nil {
		t.Fatal("not sycning", err)
	}
	nodeInfo := make(map[string]interface{})
	if err := server.rpc.Call(&nodeInfo, "admin_nodeInfo"); err != nil {
		t.Fatal("nodeInfo:", err)
	}
	enode := nodeInfo["enode"].(string)
	t.Log("enode", enode)

	// Simple client
	client, err := startGeth("./datadirs/goerli/light", false, "--goerli", "--syncmode=light", "--nodiscover")
	defer client.kill()
	if err != nil {
		t.Fatal(client.cmd.String(), err)
	}
	if err := client.addPeer(enode); err != nil {
		t.Fatal("addPeer", err)
	}

	var peers []interface{}
	if err := client.rpc.Call(&peers, "admin_peers"); err != nil {
		t.Fatal("peers", err)
	}
	if len(peers) == 0 {
		t.Log("Expected: # of client peers > 0")
		t.Fail()
	}
	t.Log("peers", peers)
	/*
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
	*/
}
