
package main

import (
	"flag"
	"fmt"
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
	rr := httptest.NewRecorder()
	rc := newRecaptchaCheck("recaptcha_v2_test_public.txt", "recaptcha_v2_test_secret.txt")
	handler := makeRootHandler(*wsaddr, "faucet.html", *rc)
	handler.ServeHTTP(rr, req)	
	fmt.Println("Result", rr.Result().Header)
	if rr.Code != http.StatusOK {
		t.Errorf("First expected 200 %v", rr.Code)
	}
	fmt.Println("hej", rr)

	n := 0
	for i := 0; i < 1; i++ {
		handler.ServeHTTP(rr, req)	
		if rr.Code != http.StatusOK {
			n += 1
			//t.Errorf("Second expected NOT 200 %v", rr.Code)
		}
	}
	if n == 0 {
		t.Errorf("Second expected NOT 200 %v", rr.Code)
	}
}