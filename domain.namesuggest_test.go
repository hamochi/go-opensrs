package opensrs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// Example 1
// For lookup, suggestion, premium, and personal names with suggestion
// limited to .COM, .NET, and .ORG, in English, German, Italian, and Spanish.
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

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "5"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "taken"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "searchstring.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "13"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "searchstring.net"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)

	}

	// Premium Items[0].Price
	if want := "5499"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "96"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "bestsearchstring.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)
	}

	// PersonalNames IsSuccess
	if want := true; resp.Attributes.PersonalNames.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.PersonalNames.IsSuccess, want %v, got %v", want, resp.Attributes.PersonalNames.IsSuccess)
	}

	// PersonalNames Count
	if want := "0"; resp.Attributes.PersonalNames.Count != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Count, want %s, got %s", want, resp.Attributes.PersonalNames.Count)
	}

}

// Example 2
// For both lookup and suggestion with lookups limited to .COM and .INFO,
// querying the registry (not OpenSRS cache) for .COM lookups, suggestions
// limited to .COM and .ORG, and maximum 25 suggestions returned.
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
				TLDs:       []string{".com", ".info"},
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

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "8"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "taken"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "examplesearch.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "22"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "examplefind.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)
	}

}

// Example 3
// For premium domains, limited to .COM and .NET.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-3
func TestNameSuggestExample3(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example3.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="searchstring">abc&amp;amp;d !</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc></dt_assoc></item><item key="premium"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.net</item></dt_array></item></dt_assoc></item><item key="suggestion"><dt_assoc></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">premium</item></dt_array></item></dt_assoc></item><item key="action">NAME_SUGGEST</item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "abc&amp;d !",
		Services:     []string{"premium"},
		ServiceOverride: NameSuggestServiceOverride{
			Premium: NameSuggestPremium{
				TLDs: []string{".com", ".net"},
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

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "4"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "abc-and-d.com"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)
	}

	// Premium Items[0].Price
	if want := "299.98"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

}

// Example 4
// For premium, lookup, and suggestion, limited to .COM, maximum 10
// suggestions returned.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-4
func TestNameSuggestExample4(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example4.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="searchstring">abc&amp;amp;d</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item><item key="premium"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item><item key="suggestion"><dt_assoc><item key="maximum">10</item><item key="tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">lookup</item><item key="1">suggestion</item><item key="2">premium</item></dt_array></item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "abc&amp;d",
		Services:     []string{"lookup", "suggestion", "premium"},
		ServiceOverride: NameSuggestServiceOverride{
			Premium: NameSuggestPremium{
				TLDs: []string{".com"},
			},
			Suggestion: NameSuggestSuggestion{
				TLDs:    []string{".com"},
				Maximum: "10",
			},
			Lookup: NameSuggestLookup{
				TLDs: []string{".com"},
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

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "2"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "available"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "abc-d.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)

	}

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "3"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "abc-and-d.com"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)

	}

	// Premium Items[0].Price
	if want := "299.98"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "10"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "abcdlive.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)

	}

}

// Example 5
// For lookup and suggestion, limited to .COM, not checking the availability of
// the lookup domain, maximum 10 suggestions returned.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-5
func TestNameSuggestExample5(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example5.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="services"><dt_array><item key="0">lookup</item><item key="1">suggestion</item></dt_array></item><item key="searchstring">smith</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item><item key="premium"><dt_assoc></dt_assoc></item><item key="suggestion"><dt_assoc><item key="maximum">10</item><item key="tlds"><dt_array><item key="0">.com</item></dt_array></item></dt_assoc></item></dt_assoc></item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "smith",
		Services:     []string{"lookup", "suggestion"},
		ServiceOverride: NameSuggestServiceOverride{
			Suggestion: NameSuggestSuggestion{
				TLDs:    []string{".com"},
				Maximum: "10",
			},
			Lookup: NameSuggestLookup{
				TLDs: []string{".com"},
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

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "1"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "undetermined"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "smith.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)

	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "10"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "myjosephsmith.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)

	}

}

// Example 6
// For lookup, suggestion, premium, and personal names with suggestion
// limited to .COM, .NET, .ORG, and .IN, in English, German, Italian, and
// Spanish, with the command run time limited to 0.4 seconds.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-6
func TestNameSuggestExample6(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example6.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.net</item><item key="2">.org</item><item key="3">in</item></dt_array></item><item key="languages"><dt_array><item key="0">en</item><item key="1">de</item><item key="2">it</item><item key="3">es</item></dt_array></item><item key="max_wait_time">0.4</item><item key="searchstring">search string</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc></dt_assoc></item><item key="premium"><dt_assoc></dt_assoc></item><item key="suggestion"><dt_assoc></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">lookup</item><item key="1">suggestion</item><item key="2">premium</item><item key="3">personal_names</item></dt_array></item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "search string",
		Services:     []string{"lookup", "suggestion", "premium", "personal_names"},
		MaxWaitTime:  "0.4",
		Languages:    []string{"en", "de", "it", "es"},
		TLDs:         []string{".com", ".net", ".org", "in"},
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "8"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "taken"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "searchstring.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)

	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "48"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "mysearchstring.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)

	}

	// PersonalNames IsSuccess
	if want := true; resp.Attributes.PersonalNames.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.PersonalNames.IsSuccess, want %v, got %v", want, resp.Attributes.PersonalNames.IsSuccess)
	}

	// PersonalNames Count
	if want := "5"; resp.Attributes.PersonalNames.Count != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Count, want %s, got %s", want, resp.Attributes.PersonalNames.Count)
	}

	// PersonalNames Item[0].Status
	if want := "available"; resp.Attributes.PersonalNames.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Items[0].Status, want %s, got %s", want, resp.Attributes.PersonalNames.Items[0].Status)
	}

	// PersonalNames Item[0].Domain
	if want := "search.stringham.com"; resp.Attributes.PersonalNames.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Items[0].Domain, want %s, got %s", want, resp.Attributes.PersonalNames.Items[0].Domain)

	}

}

// Example 7
// Resubmits the previously run name_suggest command which did not return
// complete lookup results during the specified max_wait_time. The command
// can run for a maximum of 0.7 seconds.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-7
func TestNameSuggestExample7(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example7.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version="1.0" encoding="UTF-8" standalone="no"?><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="attributes"><dt_assoc><item key="search_key">vgL2FeBzZ8JuS5lIluIEYhDc7Vg</item><item key="max_wait_time">0.7</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchKey:   "vgL2FeBzZ8JuS5lIluIEYhDc7Vg",
		MaxWaitTime: "0.7",
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "84"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "taken"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "searchstring.com"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "20"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "badmintonstring.com"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)

	}

	// Premium Items[0].Price
	if want := "1349.00"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

	// Suggestion IsSuccess
	if want := true; resp.Attributes.Suggestion.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Suggestion.IsSuccess, want %v, got %v", want, resp.Attributes.Suggestion.IsSuccess)
	}

	// Suggestion Count
	if want := "50"; resp.Attributes.Suggestion.Count != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Count, want %s, got %s", want, resp.Attributes.Suggestion.Count)
	}

	// Suggestion Item[0].Status
	if want := "available"; resp.Attributes.Suggestion.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Status, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Status)
	}

	// Suggestion Item[0].Domain
	if want := "amazonsearchstring.com"; resp.Attributes.Suggestion.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Suggestion.Items[0].Domain, want %s, got %s", want, resp.Attributes.Suggestion.Items[0].Domain)

	}

	// PersonalNames IsSuccess
	if want := true; resp.Attributes.PersonalNames.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.PersonalNames.IsSuccess, want %v, got %v", want, resp.Attributes.PersonalNames.IsSuccess)
	}

	// PersonalNames Count
	if want := "6"; resp.Attributes.PersonalNames.Count != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Count, want %s, got %s", want, resp.Attributes.PersonalNames.Count)
	}

	// PersonalNames Item[0].Status
	if want := "available"; resp.Attributes.PersonalNames.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Items[0].Status, want %s, got %s", want, resp.Attributes.PersonalNames.Items[0].Status)
	}

	// PersonalNames Item[0].Domain
	if want := "search.stringham.com"; resp.Attributes.PersonalNames.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.PersonalNames.Items[0].Domain, want %s, got %s", want, resp.Attributes.PersonalNames.Items[0].Domain)

	}

}

// Example 8
// For premium domains, limited to .COM and .NET names that cost between
// $100 and $10000, maximum 10 suggestions returned.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-8
func TestNameSuggestExample8(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example8.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="attributes"><dt_assoc><item key="searchstring">computerstore</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc></dt_assoc></item><item key="premium"><dt_assoc><item key="maximum">10</item><item key="price_max">10000</item><item key="price_min">100</item><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.net</item></dt_array></item></dt_assoc></item><item key="suggestion"><dt_assoc></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">premium</item></dt_array></item></dt_assoc></item><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item><item key="protocol">XCP</item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "computerstore",
		Services:     []string{"premium"},
		ServiceOverride: NameSuggestServiceOverride{
			Premium: NameSuggestPremium{
				TLDs:     []string{".com", ".net"},
				Maximum:  "10",
				PriceMin: "100",
				PriceMax: "10000",
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

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "2"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "childhoodneglect.com"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)

	}

	// Premium Items[0].Price
	if want := "499.00"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

}

// Example 9
// Lookup for one of the new TLDs.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-9
func TestNameSuggestExample9(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example9.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="object">DOMAIN</item><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="services"><dt_array><item key="0">lookup</item></dt_array></item><item key="tlds"><dt_array><item key="0">guru</item></dt_array></item><item key="searchstring">example</item><item key="service_override"><dt_assoc><item key="suggestion"><dt_assoc></dt_assoc></item><item key="lookup"><dt_assoc></dt_assoc></item><item key="premium"><dt_assoc></dt_assoc></item></dt_assoc></item></dt_assoc></item><item key="action">NAME_SUGGEST</item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "example",
		Services:     []string{"lookup"},
		TLDs:         []string{"guru"},
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "1"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "available"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "example.guru"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].HasClaim
	if want := true; resp.Attributes.Lookup.Items[0].HasClaim != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].HasClaim, want %v, got %v", want, resp.Attributes.Lookup.Items[0].HasClaim)
	}

}

// Example 10
// For premium, premium_brokered_transfer, and premium_make_offer
// domains, limited to .COM, .NET, .ORG, and .DE names.
// https://domains.opensrs.guide/docs/name_suggest-domain-1#example-10
func TestNameSuggestExample10(t *testing.T) {
	setup()
	defer teardown()

	respXML := readFile(t, "testresponses/domain.namesuggest.example10.xml")

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		// Test request body
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="attributes"><dt_assoc><item key="searchstring">testdomain</item><item key="service_override"><dt_assoc><item key="lookup"><dt_assoc></dt_assoc></item><item key="premium"><dt_assoc></dt_assoc></item><item key="suggestion"><dt_assoc></dt_assoc></item></dt_assoc></item><item key="services"><dt_array><item key="0">premium</item><item key="1">premium_make_offer</item><item key="2">premium_brokered_transfer</item><item key="3">lookup</item></dt_array></item><item key="tlds"><dt_array><item key="0">.com</item><item key="1">.net</item><item key="2">.org</item><item key="3">.de</item></dt_array></item></dt_assoc></item><item key="action">NAME_SUGGEST</item><item key="object">DOMAIN</item></dt_assoc></data_block></body></OPS_envelope>`
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
		SearchString: "testdomain",
		Services:     []string{"premium", "premium_make_offer", "premium_brokered_transfer", "lookup"},
		TLDs:         []string{".com", ".net", ".org", ".de"},
	})

	if err != nil {
		t.Error("unexpected error", err.Error())
		return
	}

	if resp.IsSuccess != true {
		t.Errorf("unexpected IsSuccess, want true, got %v", resp.IsSuccess)
	}

	// Lookup IsSuccess
	if want := true; resp.Attributes.Lookup.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Lookup.IsSuccess, want %v, got %v", want, resp.Attributes.Lookup.IsSuccess)
	}

	// Lookup Count
	if want := "3"; resp.Attributes.Lookup.Count != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Count, want %s, got %s", want, resp.Attributes.Lookup.Count)
	}

	// Lookup Items[0].Status
	if want := "taken"; resp.Attributes.Lookup.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Status, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Lookup Items[0].Domain
	if want := "testdomain.net"; resp.Attributes.Lookup.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Lookup.Items[0].Domain, want %s, got %s", want, resp.Attributes.Lookup.Items[0].Status)
	}

	// Premium IsSuccess
	if want := true; resp.Attributes.Premium.IsSuccess != Bool(want) {
		t.Errorf("unexpected resp.Attributes.Premium.IsSuccess, want %v, got %v", want, resp.Attributes.Premium.IsSuccess)
	}

	// Premium Count
	if want := "15"; resp.Attributes.Premium.Count != want {
		t.Errorf("unexpected resp.Attributes.Premium.Count, want %s, got %s", want, resp.Attributes.Premium.Count)
	}

	// Premium Items[0].Status
	if want := "available"; resp.Attributes.Premium.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Status, want %s, got %s", want, resp.Attributes.Premium.Items[0].Status)
	}

	// Premium Items[0].Domain
	if want := "testdomain.com"; resp.Attributes.Premium.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Domain, want %s, got %s", want, resp.Attributes.Premium.Items[0].Domain)

	}

	// Premium Items[0].Price
	if want := "6999"; resp.Attributes.Premium.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.Premium.Items[0].Price, want %s, got %s", want, resp.Attributes.Premium.Items[0].Price)
	}

	// PremiumBrokeredTransfer Count
	if want := "2"; resp.Attributes.PremiumBrokeredTransfer.Count != want {
		t.Errorf("unexpected resp.Attributes.PremiumBrokeredTransfer.Count, want %s, got %s", want, resp.Attributes.PremiumBrokeredTransfer.Count)
	}

	// PremiumBrokeredTransfer Items[0].Status
	if want := "available"; resp.Attributes.PremiumBrokeredTransfer.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.PremiumBrokeredTransfer.Items[0].Status, want %s, got %s", want, resp.Attributes.PremiumBrokeredTransfer.Items[0].Status)
	}

	// PremiumBrokeredTransfer Items[0].Domain
	if want := "wangtestdomain.com"; resp.Attributes.PremiumBrokeredTransfer.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.PremiumBrokeredTransfer.Items[0].Domain, want %s, got %s", want, resp.Attributes.PremiumBrokeredTransfer.Items[0].Domain)
	}

	// PremiumBrokeredTransfer Items[0].Price
	if want := "1147"; resp.Attributes.PremiumBrokeredTransfer.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.PremiumBrokeredTransfer.Items[0].Price, want %s, got %s", want, resp.Attributes.PremiumBrokeredTransfer.Items[0].Price)
	}

	// PremiumMakeOffer Count
	if want := "3"; resp.Attributes.PremiumMakeOffer.Count != want {
		t.Errorf("unexpected resp.Attributes.PremiumMakeOffer.Count, want %s, got %s", want, resp.Attributes.PremiumMakeOffer.Count)
	}

	// PremiumMakeOffer Items[0].Status
	if want := "available"; resp.Attributes.PremiumMakeOffer.Items[0].Status != want {
		t.Errorf("unexpected resp.Attributes.PremiumMakeOffer.Items[0].Status, want %s, got %s", want, resp.Attributes.PremiumMakeOffer.Items[0].Status)
	}

	// PremiumMakeOffer Items[0].Domain
	if want := "testdomain.net"; resp.Attributes.PremiumMakeOffer.Items[0].Domain != want {
		t.Errorf("unexpected resp.Attributes.PremiumMakeOffer.Items[0].Domain, want %s, got %s", want, resp.Attributes.PremiumMakeOffer.Items[0].Domain)
	}

	// PremiumMakeOffer Items[0].Price
	if want := "0"; resp.Attributes.PremiumMakeOffer.Items[0].Price != want {
		t.Errorf("unexpected resp.Attributes.PremiumMakeOffer.Items[0].Price, want %s, got %s", want, resp.Attributes.PremiumMakeOffer.Items[0].Price)
	}

	// PremiumMakeOffer Items[0].ThirdPartyOfferUrl
	if want := "http://www.sedo.com/search/details.php4?language=us&partnerid=316601&domain="; resp.Attributes.PremiumMakeOffer.Items[0].ThirdPartyOfferUrl != URL(want) {
		t.Errorf("unexpected resp.Attributes.PremiumMakeOffer.Items[0].Price, \nwant %s,\ngot %s", want, resp.Attributes.PremiumMakeOffer.Items[0].ThirdPartyOfferUrl)
	}

}
