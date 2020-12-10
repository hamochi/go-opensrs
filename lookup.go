package opensrs

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type LookupResponse struct {
	Available bool
}

func (c *Client) Lookup(domainName string, noCache bool) (LookupResponse, error) {
	lResp := LookupResponse{}

	noCacheValue := "0"
	if noCache == true {
		noCacheValue = "1"
	}

	resp, statusCode, err := c.do("DOMAIN", "LOOKUP",
		item{Key: "domain", Value: domainName},
		item{Key: "no_cache", Value: noCacheValue})

	if err != nil {
		return lResp, err
	}

	if statusCode != http.StatusOK {
		return lResp, errors.New(fmt.Sprintf("unexpected status code from Lookup, expected 200 got %d", statusCode))
	}

	var respXml ApiRequestResponse
	err = xml.Unmarshal(resp, &respXml)
	if err != nil {
		return lResp, err
	}

	for _, item := range respXml.DtAssoc.Items {
		if item.Key == "attributes" {
			dtAssoc := *item.DtAssoc
			for _, attrItem := range dtAssoc.Items {
				if attrItem.Key == "status" {
					if attrItem.Value == "available" {
						lResp.Available = true
						return lResp, nil
					}
					return lResp, nil
				}
			}
		}
	}

	return lResp, errors.New(fmt.Sprintf("unexpected response, resp:%s", string(resp)))
}
