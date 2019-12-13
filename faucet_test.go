
package main

import (
	"net/http"
    "testing"
)

func TestThrottling(t *testing.T) {
	resp1, err := http.Get("http://127.0.0.1:8088")
	if err != nil {
		t.Errorf("first %s", err)
	}
	defer resp1.Body.Close()
	resp2, err := http.Get("http://127.0.0.1:8088")
	if resp2.StatusCode == http.StatusOK {
		t.Errorf("Expected error for second %s", resp2.Status)
	}
	defer resp2.Body.Close()
}