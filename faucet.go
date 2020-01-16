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
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/ezzarghili/recaptcha-go.v3"
)

type balanceInfo struct {
	BalanceBefore int
	BalanceAfter  int
}

func addBalance(serverRPCEndpoint string, clientNodeID string, balance int, topic string) (balanceInfo, error) {
	var info balanceInfo
	server, err := rpc.Dial(serverRPCEndpoint)
	if err != nil {
		return info, fmt.Errorf("rpc.Dial %s: %s", serverRPCEndpoint, err)
	}
	log.Printf("Server connected")

	var balances []int
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.CallContext(ctx, &balances, "les_addBalance", clientNodeID, balance, topic); err != nil {
		return info, fmt.Errorf("les_addBalance %s", err)
	}
	info.BalanceBefore = balances[0]
	info.BalanceAfter = balances[1]
	log.Println("AddBalance success", balances)

	return info, nil
}

type allClientsInfo map[string]map[string]interface{}
type clientInfo map[string]interface{}

func getBalance(serverRPCEndpoint string, nodeID string) (clientInfo, error) {
	var allInfo allClientsInfo
	var cInfo clientInfo
	log.Printf("GetBalance called at %s for %s", serverRPCEndpoint, nodeID)
	server, err := rpc.Dial(serverRPCEndpoint)
	if err != nil {
		return cInfo, fmt.Errorf("rpc.Dial %s: %s", serverRPCEndpoint, err)
	}
	log.Printf("Server connected")

	clientIDs := []string{nodeID}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := server.CallContext(ctx, &allInfo, "les_clientInfo", clientIDs); err != nil {
		return cInfo, fmt.Errorf("les_clientinfo: %s", err)
	}
	cInfo = allInfo[nodeID]
	log.Printf("Got balance %s", allInfo)
	return cInfo, nil
}

type formData struct {
	NodeID      string
	Error       error
	Balance     balanceInfo
	Client      clientInfo
	RPCEndpoint string
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
		log.Println("handle", r.Method, r.Form)

		switch {
		case nodeID == "":
			err = errors.New("nodeID is required")
		case r.Method == http.MethodPost:
			err = rc.checker.Verify(r.FormValue("g-recaptcha-response"))
			if err != nil {
				break
			}
			bInfo, err = addBalance(rpcEndpoint, nodeID, 1000, "foobar")
			if err != nil {
				break
			}
			cInfo, err = getBalance(rpcEndpoint, nodeID)
			if cInfo != nil {
				cInfo["pricing/oldBalance"] = bInfo.BalanceBefore
			}
		case r.Method == http.MethodGet:
			cInfo, err = getBalance(rpcEndpoint, nodeID)
		default:
			err = errors.New("Unsupported method")
		}

		fillData := formData{nodeID, err, bInfo, cInfo, rpcEndpoint, rc.public}

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
