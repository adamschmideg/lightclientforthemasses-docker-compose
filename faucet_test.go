
package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"
)


var wsaddr = flag.String("wsaddr", "127.0.0.1", "Address of the faucet")
func Init() {
	flag.Parse()
}

func TestThrottling(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Real-IP", "::1")
	rr := httptest.NewRecorder()
	rc := newRecaptchaCheck("recaptcha_v2_test_public.txt", "recaptcha_v2_test_secret.txt")
	handler := makeRootHandler(*wsaddr, "faucet.html", *rc)

	handler.ServeHTTP(rr, req)	
	if rr.Code != http.StatusOK {
		t.Errorf("First expected 200 %v", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)	
	if rr.Code == http.StatusOK {
		t.Errorf("Second expected NOT 200 %v", rr.Code)
	}
}