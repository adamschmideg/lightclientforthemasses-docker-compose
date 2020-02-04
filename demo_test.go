package main

import (
	"os/exec"
	"fmt"
	"log"
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

var rpcArgs = []string{"--rpc", "--rpcapi=admin,eth,les", "--rpcaddr=0.0.0.0"}

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
	allArgs := []string{"--datadir", datadir}
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
	allArgs := []string{"--datadir", g.datadir}
	execArgs := []string{"attach", "--exec", js}
	allArgs = append(allArgs, rpcArgs...)
	allArgs = append(allArgs, execArgs...)
	cmd := exec.Command("geth", allArgs...)
	log.Println("to exec", cmd.String())
	var b []byte
	b, err := cmd.CombinedOutput()
	out := strings.Trim(string(b), " \n\r\t\"")
	if err != nil {
		return out, err
	}
	return out, nil
}

func (g *geth) kill() error {
	return g.cmd.Process.Kill()
}

func TestDemo(t *testing.T) {
	server, err := startGeth("./datadirs/goerli/fast", true, "--lightserv=100", "--light.maxpeers=1", "--goerli", "--syncmode=fast")
	if err != nil {
		t.Error(server.cmd.String(), err)
	}
	enode, err := server.exec("admin.nodeInfo.enode")
	if err != nil {
		t.Error(enode, err)
	}
	client, err := startGeth("./datadirs/goerli/light", false, "--syncmode=light", "--nodiscover")
	if err != nil {
		t.Error(server.cmd.String(), err)
	}
	addPeerJs := fmt.Sprintf(`'admin.addPeer("%v")'`, enode)
	addPeerResult, err := server.exec(addPeerJs)
	if err != nil {
		t.Error(addPeerResult, err)
	}
	t.Log(enode)
	t.Log("addPeer", addPeerResult)
	server.kill()
	client.kill()
}

// Start a priority client with an empty datadir
// Get the nodeID of the priority client
// Add balance for the priority client on the light server
// Add the server as peer to let priority client start syncing
// Check if it's actually syncing
// Check if the regular client got kicked out
