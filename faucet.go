package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/ezzarghili/recaptcha-go.v3"
)

// How many tokens the faucet gives at a request
const tokensPerMinit = 3_000_000_000

type internalTokens int
type minits int

func toMinits(n internalTokens) minits {
	return minits(n / tokensPerMinit)
}

func fromMinits(n minits) internalTokens {
	return internalTokens(n * tokensPerMinit)
}

type balanceInfo struct {
	BalanceBefore minits
	BalanceAfter  minits
}

type client struct { 
	api *rpc.Client
}

func newClient(rpcEndpoint string) *client {
	c := &client{}
	var err error
	c.api, err = rpc.Dial(rpcEndpoint)
	if err != nil {
		fmt.Errorf("rpc.Dial %s: %s", rpcEndpoint, err)
	}
	return c
}

func (c* client) addBalance(clientNodeID string, minitsToAdd minits, topic string) (balanceInfo,error) {
	var info balanceInfo

	var balances []internalTokens
	tokens := fromMinits(minitsToAdd)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.api.CallContext(ctx, &balances, "les_addBalance", clientNodeID, tokens, topic); err != nil {
		return info, fmt.Errorf("les_addBalance %s", err)
	}
	info.BalanceBefore = toMinits(balances[0])
	info.BalanceAfter = toMinits(balances[1])
	log.Println("AddBalance success", balances)

	return info, nil
}

type allClientsInfo map[string]map[string]interface{}
type clientInfo map[string]interface{}

func (c* client) getBalance(nodeID string) (clientInfo,error) {
	var allInfo allClientsInfo
	var cInfo clientInfo
	log.Printf("GetBalance called for %s", nodeID)

	clientIDs := []string{nodeID}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.api.CallContext(ctx, &allInfo, "les_clientInfo", clientIDs); err != nil {
		return cInfo, fmt.Errorf("les_clientinfo: %s", err)
	}
	cInfo = allInfo[nodeID]
	log.Printf("Got balance %s", allInfo)
	return cInfo, nil
}

func (c* client) getServerENode() (string,error) {
	var nodeInfo p2p.NodeInfo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := c.api.CallContext(ctx, &nodeInfo, "admin_nodeInfo"); err != nil {
		return "", fmt.Errorf("admin_nodeInfo: %s", err)
	}
	log.Printf("Got nodeInfo %s", nodeInfo.Enode)
	return nodeInfo.Enode, nil

}

type formData struct {
	NodeID      string
	Error       error
	Balance     balanceInfo
	Client      clientInfo
	ServerENode string
	Recaptcha   string
}

type recaptchaCheck struct {
	checker recaptcha.ReCAPTCHA
	public  string
}

func newRecaptchaCheck(publicPath string, secretPath string) *recaptchaCheck {
	rc := recaptchaCheck{}
	var err error
	data, err := ioutil.ReadFile(publicPath)
	if err != nil {
		fmt.Println("Can't open recaptcha public", err)
		return &rc
	}
	rc.public = string(data)
	data, err = ioutil.ReadFile(secretPath)
	if err != nil {
		fmt.Println("Can't open recaptcha secret", err)
		return &rc
	}
	recaptchaSecret := string(data)
	recaptchaChecker, err := recaptcha.NewReCAPTCHA(recaptchaSecret, recaptcha.V2, 10*time.Second)
	if err != nil {
		fmt.Println("Can't make recaptchaChecker")
		return &rc
	}
	rc.checker = recaptchaChecker
	return &rc
}

func makeRootHandler(rpcEndpoint string, templatePath string, rc recaptchaCheck) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		nodeID := r.FormValue("nodeID")
		var err error
		var bInfo balanceInfo
		var cInfo clientInfo
		var serverENode string
		log.Println("handle", r.Method, r.Form)
		c := newClient(rpcEndpoint)

		switch {
		case nodeID == "":
			err = errors.New("nodeID is required")
		case r.Method == http.MethodPost:
			err = rc.checker.Verify(r.FormValue("g-recaptcha-response"))
			if err != nil {
				break
			}
			bInfo, err = c.addBalance(nodeID, 1, "foobar")
			if err != nil {
				break
			}
			cInfo, err = c.getBalance(nodeID)
			if cInfo != nil {
				cInfo["pricing/oldBalance"] = bInfo.BalanceBefore
			}
		case r.Method == http.MethodGet:
			cInfo, err = c.getBalance(nodeID)
			serverENode, err = c.getServerENode()
		default:
			err = errors.New("Unsupported method")
		}

		fillData := formData{nodeID, err, bInfo, cInfo, serverENode, rc.public}

		t, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Println("Parsing html", err)
			fmt.Fprintln(w, "Internal error")
			return
		}
		// fmt.Println("fillData", fillData)
		if err := t.Execute(w, fillData); err != nil {
			log.Println("Executing template", err)
			fmt.Fprintln(w, "internal error")
		}
	}

	lmt := rateLimiter()
	handler := http.HandlerFunc(handlerFunc)
	return tollbooth.LimitHandler(lmt, handler)
}

func lookupIP(address string) string {
	ips, _ := net.LookupIP(address)
	if len(ips) > 0 {
		return ips[0].String()
	}
	return address
}

func rateLimiter() *limiter.Limiter {
	lmt := tollbooth.NewLimiter(1, nil)
	lmt.SetStatusCode(http.StatusTooManyRequests)
	return lmt
}

func main() {
	rpcaddr := flag.String("rpcaddr", "127.0.0.1", "Address of the lightserver's rpc endpoint")
	rpcport := flag.Int("rpcport", 8545, "Port of the lightserver's rpc endpoint")
	port := flag.Int("port", 8088, "Web service port of the faucet")
	templatePath := flag.String("template", "/var/www/faucet.html", "Full path to the html template file")
	recaptchaPublicFile := flag.String("recaptcha.public", "recaptcha_v2_test_public.txt", "Path to public key")
	recaptchaSecretFile := flag.String("recaptcha.secret", "recaptcha_v2_test_secret.txt", "Path to secret key")
	flag.Parse()

	recaptchaCheck := newRecaptchaCheck(*recaptchaPublicFile, *recaptchaSecretFile)
	// I have to resolve to IP address inside a docker container, it's not working with a name
	rpcIP := lookupIP(*rpcaddr)
	rpcEndpoint := fmt.Sprintf("http://%s:%v", rpcIP, *rpcport)
	rootHandler := makeRootHandler(rpcEndpoint, *templatePath, *recaptchaCheck)
	wsAddress := fmt.Sprintf(":%v", *port)
	log.Println("Listening at", wsAddress, ", calling", rpcEndpoint)

	http.Handle("/", rootHandler)
	log.Fatal(http.ListenAndServe(wsAddress, nil))
}
