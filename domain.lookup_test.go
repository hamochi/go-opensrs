package opensrs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// https://domains.opensrs.guide/docs/lookup-domain-2#example-1
func TestLookupExample1a(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.lookup.example1a.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">LOOKUP</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="domain">example.com</item><item key="no_cache">1</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		var wantReqXml LookupRequest
		err = FromXml([]byte(want), &wantReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"wanted request\" : ", err.Error())
		}

		var gotReqXml LookupRequest
		err = FromXml(body, &gotReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"got request\":  ", err.Error())
		}

		if !reflect.DeepEqual(wantReqXml, gotReqXml) {
			t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", gotReqXml, wantReqXml)
		}

		// Test req method
		testMethod(t, r)

		// Test authentication method
		testAuth(t, r.Header, string(body))

		fmt.Fprint(w, respXML)
	})

	resp, err := client.Domains.Lookup(LookupRequestAttributes{
		Domain:  "example.com",
		NoCache: true,
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	want := &LookupResponse{
		BaseResponse: BaseResponse{
			Action:       "REPLY",
			Object:       "DOMAIN",
			Protocol:     "XCP",
			ResponseCode: "210",
			IsSuccess:    true,
			ResponseText: "Domain available",
		},
		Attributes: LookupResponseAttributes{
			Status: "available",
		},
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", resp, want)
	}

}

// https://domains.opensrs.guide/docs/lookup-domain-2#example-1
func TestLookupExample1b(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.lookup.example1b.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">LOOKUP</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="domain">example.com</item><item key="no_cache">1</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		var wantReqXml LookupRequest
		err = FromXml([]byte(want), &wantReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"wanted request\" : ", err.Error())
		}

		var gotReqXml LookupRequest
		err = FromXml(body, &gotReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"got request\":  ", err.Error())
		}

		if !reflect.DeepEqual(wantReqXml, gotReqXml) {
			t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", gotReqXml, wantReqXml)
		}

		// Test req method
		testMethod(t, r)

		// Test authentication method
		testAuth(t, r.Header, string(body))

		fmt.Fprint(w, respXML)
	})

	resp, err := client.Domains.Lookup(LookupRequestAttributes{
		Domain:  "example.com",
		NoCache: true,
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	want := &LookupResponse{
		BaseResponse: BaseResponse{
			Action:       "REPLY",
			Object:       "DOMAIN",
			Protocol:     "XCP",
			ResponseCode: "211",
			IsSuccess:    true,
			ResponseText: "Domain taken",
		},
		Attributes: LookupResponseAttributes{
			Status: "taken",
		},
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", resp, want)
	}

}

// https://domains.opensrs.guide/docs/lookup-domain-2#example-2
func TestLookupExample2(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.lookup.example2.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">LOOKUP</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="domain">example.guru</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		var wantReqXml LookupRequest
		err = FromXml([]byte(want), &wantReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"wanted request\" : ", err.Error())
		}

		var gotReqXml LookupRequest
		err = FromXml(body, &gotReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"got request\":  ", err.Error())
		}

		if !reflect.DeepEqual(wantReqXml, gotReqXml) {
			t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", gotReqXml, wantReqXml)
		}

		// Test req method
		testMethod(t, r)

		// Test authentication method
		testAuth(t, r.Header, string(body))

		fmt.Fprint(w, respXML)
	})

	resp, err := client.Domains.Lookup(LookupRequestAttributes{
		Domain: "example.guru",
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	want := &LookupResponse{
		BaseResponse: BaseResponse{
			Action:       "REPLY",
			Object:       "DOMAIN",
			Protocol:     "XCP",
			ResponseCode: "210",
			IsSuccess:    true,
			ResponseText: "Domain available",
		},
		Attributes: LookupResponseAttributes{
			Status:   "available",
			HasClaim: true,
		},
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("lookup returned, got\n%+v,\nwant\n%+v", resp, want)
	}

}
