package opensrs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("testUser", "testApi")
	client.BaseURL = server.URL + "/"
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request) {
	if r.Method != http.MethodPost {
		t.Errorf("Want POST, got %s", r.Method)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	if string(b) != want {
		t.Errorf(fmt.Sprintf("Unexpected body, want \n%s\n got\n%s", want, string(b)))
	}
}

func TestBuildXMLRequest(t *testing.T) {
	setup()
	defer teardown()

	xml, err := client.buildXMLRequest("DOMAIN", "LOOKUP", item{Key: "domain", Value: "test.com"})
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	fmt.Println(xml)
	expected := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="object">DOMAIN</item><item key="action">LOOKUP</item><item key="attributes"><dt_assoc><item key="domain">test.com</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
	if xml != expected {
		t.Errorf("Not expected this, got \n%s\nwant\n%s", xml, expected)
	}
}

func TestBuildHttpRequest(t *testing.T) {
	setup()
	defer teardown()

	client.BaseURL = "https://fake-api-server:4444"

	xmlReq := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="object">DOMAIN</item><item key="action">LOOKUP</item><item key="attributes"><dt_assoc></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`

	req, err := client.buildHttpRequest(xmlReq)
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	if req.Method != http.MethodPost {
		t.Error("incorrect http method")
	}
	if req.Header.Get("Content-Type") != "text/xml" {
		t.Error("incorrect Content-Type")
	}
	if req.Header.Get("X-Username") != "testUser" {
		t.Error("incorrect Username")
	}
	if req.Header.Get("X-Signature") != "7f9c1daee2c72416278e3188f2b62d51" {
		t.Error("incorrect Signature")
	}

	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	if string(b) != xmlReq {
		t.Errorf("Incorrect body, got \n%s\nwant\n%s", string(b), xmlReq)
	}

	if req.URL.Port() != "4444" {
		t.Errorf("incorrect Host, want 4444, got %s", req.URL.Port())
	}
	if req.URL.String() != "https://fake-api-server:4444" {
		t.Errorf("incorrect Host, want ssss, got %s", req.URL.String())
	}

}
