// Package opensrs provides a client for the OpenSRS API.
// In order to use this package you will need a OpenSRS account.
package opensrs

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Version          = "0.0.1"
	defaultBaseURL   = "https://rr-n1-tor.opensrs.net:55443"
	defaultUserAgent = "opensrs-go/" + Version
	xmlHeader        = "<?xml version='1.0' encoding='UTF-8' standalone='no' ?><!DOCTYPE OPS_envelope SYSTEM 'ops.dtd'>"
)

type BaseRequest struct {
	Action   string `json:"action"`
	Object   string `json:"object"`
	Protocol string `json:"protocol"`
}

type BaseResponse struct {
	Action       string `json:"action"`
	Object       string `json:"object"`
	Protocol     string `json:"protocol"`
	IsSuccess    Bool   `json:"is_success"`
	ResponseCode string `json:"response_code"`
	ResponseText string `json:"response_text"`
}

type Client struct {
	HttpClient       *http.Client
	ApiKey           string
	ResellerUsername string
	BaseURL          string
	Debug            bool
	Domains          *DomainsService
}

func NewClient(ResellerUsername, ApiKey string) *Client {
	c := &Client{
		ApiKey:           ApiKey,
		ResellerUsername: ResellerUsername,
		HttpClient:       &http.Client{},
		BaseURL:          defaultBaseURL,
	}
	c.Domains = &DomainsService{Client: c}
	return c
}

func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	url := c.BaseURL + path

	body := new(bytes.Buffer)
	if payload != nil {
		xml, err := ToXml(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer([]byte(xmlHeader))
		body.Write(xml)

		if c.Debug {
			log.Printf("Requst sent: %s\n", xmlHeader+string(xml))
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	signature := getMD5Hash(getMD5Hash(body.String()+c.ApiKey) + c.ApiKey)

	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("X-Username", c.ResellerUsername)
	req.Header.Set("X-Signature", signature)

	return req, nil
}

func (c *Client) Do(req *http.Request, obj interface{}) error {
	if c.Debug {
		log.Printf("Executing request (%v): %#v", req.URL, req)
	}
	e := ErrorResponse{}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		e.Err = err
		return e
	}
	defer resp.Body.Close()

	e.HttpResponse = resp

	if c.Debug {
		log.Printf("Response received: %#v\n", resp)
	}

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		e.Err = errors.New("unexpected return code")
		return e
	}

	if obj != nil {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			e.Err = err
			return e
		}
		if c.Debug {
			log.Printf("Resp body received:\n##########\n%s\n##########\n", string(b))
		}

		err = FromXml(b, obj)
		if err != nil {
			e.Err = err
			return e
		}

		oResp := BaseResponse{}
		err = FromXml(b, &oResp)
		if err != nil {
			e.Err = err
			return e
		}

		e.OpenSRSResponse = &oResp

		if oResp.IsSuccess != true {
			return e
		}

	}

	return nil
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//func Bool(b bool) *bool {
//	return &b
//}

func String(s string) *string {
	return &s
}

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	value, err := parseBool(string(data))
	if err != nil {
		return err
	}

	*b = Bool(value)
	return nil
}

//func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
//	var el *string
//	if err := d.DecodeElement(&el, &start); err != nil {
//		return err
//	}
//
//	if el == nil {
//		return nil
//	}
//
//	value, err := parseBool(*el)
//	if err != nil {
//		return err
//	}
//
//	*b = Bool(value)
//	return nil
//}

func parseBool(value string) (b bool, err error) {
	switch value {
	case `"1"`:
		b = true
	case `"0"`:
		b = false
	default:
		err = fmt.Errorf("invalid value for bool: %s", value)
	}

	return
}
