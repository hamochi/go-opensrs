package opensrs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-1
func TestNameSuggestExample1(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example1.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="languages"><dt_array><item key="0">en</item><item key="1">de</item><item key="2">it</item><item key="3">es</item></dt_array></item><item key="searchstring">search string</item><item key="services"><dt_array><item key="0">lookup</item><item key="1">suggestion</item><item key="2">premium</item><item key="3">personal_names</item></dt_array></item><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.net</item><item key="2">.org</item></dt_array></item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		var wantReqXml NameSuggestRequest
		err = FromXml([]byte(want), &wantReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"wanted request\" : ", err.Error())
		}

		var gotReqXml NameSuggestRequest
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

	// Build a request
	resp, err := client.Domains.NameSuggest(NameSuggestRequestAttributes{
		Services:     []string{"lookup", "suggestion", "premium", "personal_names"},
		SearchString: "search string",
		Languages:    []string{"en", "de", "it", "es"},
		TLDs:         []string{".com", ".net", ".org"},
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	if resp.Attributes.Suggestion.Count != "96" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want 96, got %s", resp.Attributes.Suggestion.Count)
	}

	if resp.Attributes.Suggestion.Items[0].Domain != "bestsearchstring.com" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want bestsearchstring.com, got %s", resp.Attributes.Suggestion.Items[0].Domain)
	}

	if resp.Attributes.Suggestion.Items[0].Status != "available" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want available, got %s", resp.Attributes.Suggestion.Items[0].Status)
	}

	if resp.Attributes.Lookup.Count != "5" {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want 5, got %s", resp.Attributes.Lookup.Count)
	}

	if resp.Attributes.Lookup.Items[0].Status != "taken" {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want taken, got %s", resp.Attributes.Lookup.Items[0].Status)
	}

}

// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-2
func TestNameSuggestExample2(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example2.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="attributes"><dt_assoc><item key="searchstring">example@search.com</item><item key="service_override"><dt_assoc><item key="suggestion"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.org</item></dt_array></item><item key="maximum">25</item></dt_assoc></item><item key="lookup"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.info</item></dt_array></item><item key="no_cache_tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">lookup</item><item key="1">suggestion</item></dt_array></item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		var wantReqXml NameSuggestRequest
		err = FromXml([]byte(want), &wantReqXml)
		if err != nil {
			t.Fatal("error unmarshalling \"wanted request\" : ", err.Error())
		}

		var gotReqXml NameSuggestRequest
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

	// Build a request
	resp, err := client.Domains.NameSuggest(NameSuggestRequestAttributes{
		SearchString: "example@search.com",
		Services:     []string{"lookup", "suggestion"},
		ServiceOverride: NameSuggestServiceOverride{
			Suggestion: NameSuggestSuggestion{
				TLDs:    []string{".com", ".org"},
				Maximum: "25",
			},
			Lookup: NameSuggestLookup{
				NameSuggestSuggestion: NameSuggestSuggestion{
					TLDs: []string{".com", ".info"},
				},
				NoCacheTld: []string{".com"},
			},
		},
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	if resp.Attributes.Suggestion.Count != "22" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want 96, got %s", resp.Attributes.Suggestion.Count)
	}

	if resp.Attributes.Suggestion.Items[0].Domain != "examplefind.com" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want bestsearchstring.com, got %s", resp.Attributes.Suggestion.Items[0].Domain)
	}

	if resp.Attributes.Suggestion.Items[0].Status != "available" {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want available, got %s", resp.Attributes.Suggestion.Items[0].Status)
	}

	if resp.Attributes.Lookup.Count != "8" {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want 8, got %s", resp.Attributes.Lookup.Count)
	}
	if resp.Attributes.Lookup.IsSuccess != true {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want true, got %v", resp.Attributes.Lookup.IsSuccess)
	}

	if resp.Attributes.Lookup.Items[0].Status != "taken" {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want taken, got %s", resp.Attributes.Lookup.Items[0].Status)
	}

}
