package opensrs

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestLookupTaken(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version='1.0' encoding="UTF-8" standalone="no" ?>
<!DOCTYPE OPS_envelope SYSTEM "ops.dtd">
<OPS_envelope>
	<header>
		<version>0.9</version>
  	</header>
 	<body>
  		<data_block>
			<dt_assoc>
				<item key="action">REPLY</item>
				<item key="object">DOMAIN</item>
				<item key="protocol">XCP</item>
				<item key="response_code">211</item>
				<item key="is_success">1</item>
				<item key="response_text">Domain taken</item>
				<item key="attributes">
					<dt_assoc>
						<item key="status">taken</item>
						<item key="reason"></item>
					</dt_assoc>
				</item>
			</dt_assoc>
		</data_block>
	</body>
</OPS_envelope>`

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		want := `<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="object">DOMAIN</item><item key="action">LOOKUP</item><item key="attributes"><dt_assoc><item key="domain">test.com</item><item key="no_cache">0</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`
		testMethod(t, r)
		testBody(t, r, want)
		fmt.Fprint(w, respXML)
	})

	resp, err := client.Lookup("test.com", false)

	if err != nil {
		t.Fatal("unexpected error", err)
	}

	want := LookupResponse{
		Available: false,
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Lookup returned %+v, want %+v", resp, want)
	}

}

func TestLookupAvailable(t *testing.T) {
	setup()
	defer teardown()

	respXML := `<?xml version='1.0' encoding="UTF-8" standalone="no" ?>
<!DOCTYPE OPS_envelope SYSTEM "ops.dtd">
<OPS_envelope>
    <header>
        <version>0.9</version>
    </header>
    <body>
        <data_block>
            <dt_assoc>
                <item key="object">DOMAIN</item>
                <item key="action">REPLY</item>
                <item key="response_text">Domain available</item>
                <item key="protocol">XCP</item>
                <item key="is_success">1</item>
                <item key="attributes">
                    <dt_assoc>
                        <item key="status">available</item>
                    </dt_assoc>
                </item>
                <item key="response_code">210</item>
            </dt_assoc>
        </data_block>
    </body>
</OPS_envelope>`

	domain := fmt.Sprintf("test%d.com", time.Now().Unix())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		want := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'><OPS_envelope><header><version>0.9</version></header><body><data_block><dt_assoc><item key="protocol">XCP</item><item key="object">DOMAIN</item><item key="action">LOOKUP</item><item key="attributes"><dt_assoc><item key="domain">%s</item><item key="no_cache">1</item></dt_assoc></item></dt_assoc></data_block></body></OPS_envelope>`, domain)
		testMethod(t, r)
		testBody(t, r, want)
		fmt.Fprint(w, respXML)
	})

	resp, err := client.Lookup(domain, true)

	if err != nil {
		t.Fatal("unexpected error", err)
	}

	want := LookupResponse{
		Available: true,
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Lookup returned %+v, want %+v", resp, want)
	}

}
