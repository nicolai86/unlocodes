package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	srv = httptest.NewServer(Handler())
	os.Exit(m.Run())
}

func TestInvalidRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/", srv.URL), nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Update status code%v:\n\tbody: %+v", "400", resp.StatusCode)
	}
}

func TestSuccessfulTranslate(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/?code=ESPMI&code=BBP", srv.URL), nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Update status code%v:\n\tbody: %+v", "200", resp.StatusCode)
	}

	translated := make(map[string]string)
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&translated)

	if palma, ok := translated["ESPMI"]; !ok {
		t.Fatalf("Failed to translate sample code ESPMI")
	} else if palma != "Palma de Mallorca" {
		t.Fatalf("Failed to correctly translate sample code ESPMI")
	}

	if _, ok := translated["BBP"]; ok {
		t.Fatalf("Should not include unknown locodes in response: %v", translated)
	}
}
