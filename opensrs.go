package opensrs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://rr-n1-tor.opensrs.net:55443"
const defaultTestBaseURL = "https://horizon.opensrs.net:55443"
const defaultTimeout = time.Second * 60

const version = 0.9
const xmlHeader = "<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'>"

// Client represents a client used to make calls to the OpenSRS API.
type Client struct {
	ResellerUsername string
	ApiKey           string
	HttpClient       *http.Client
	Timeout          time.Duration
	BaseURL          string
}

type ApiRequestResponse struct {
	XMLName xml.Name `xml:"OPS_envelope"`
	Version float32  `xml:"header>version"`
	DtAssoc dtAssoc  `xml:"body>data_block>dt_assoc"`
}

type dtAssoc struct {
	XMLName xml.Name `xml:"dt_assoc"`
	Items   []item   `xml:"item"`
}
type item struct {
	XMLName xml.Name `xml:"item"`
	Value   string   `xml:",chardata"`
	Key     string   `xml:"key,attr"`
	DtAssoc *dtAssoc
}

func (c *Client) buildXMLRequest(object, action string, items ...item) (string, error) {
	req := &ApiRequestResponse{
		Version: version,
		DtAssoc: dtAssoc{
			Items: []item{
				{Key: "protocol", Value: "XCP"},
				{Key: "object", Value: object},
				{Key: "action", Value: action},
				{
					Key: "attributes",
					DtAssoc: &dtAssoc{
						Items: items,
					},
				},
			},
		},
	}

	xmlReq, err := xml.Marshal(req)
	//xmlReq, err := xml.MarshalIndent(req, "", "   ")
	if err != nil {
		return "", err
	}

	return xmlHeader + string(xmlReq), nil
}

func (c *Client) buildHttpRequest(xmlReq string) (*http.Request, error) {
	c.HttpClient.Timeout = c.Timeout
	req, err := http.NewRequest(http.MethodPost, c.BaseURL, strings.NewReader(xmlReq))
	if err != nil {
		return nil, err
	}

	// Signature according OpenSRS docs
	signature := getMD5Hash(getMD5Hash(xmlReq+c.ApiKey) + c.ApiKey)

	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("X-Username", c.ResellerUsername)
	req.Header.Set("X-Signature", signature)

	return req, nil
}

func (c *Client) sendRequest(req *http.Request) ([]byte, int, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return buf, resp.StatusCode, nil
}

func (c *Client) do(object, action string, items ...item) ([]byte, int, error) {
	xmlReq, err := c.buildXMLRequest(object, action, items...)
	if err != nil {
		return nil, 0, err
	}

	req, err := c.buildHttpRequest(xmlReq)
	if err != nil {
		return nil, 0, nil
	}

	return c.sendRequest(req)

}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func NewClient(ResellerUsername, ApiKey string) *Client {
	return &Client{
		ResellerUsername: ResellerUsername,
		ApiKey:           ApiKey,
		HttpClient:       http.DefaultClient,
		Timeout:          defaultTimeout,
		BaseURL:          defaultBaseURL,
	}
}
