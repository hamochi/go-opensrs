package opensrs

type LookupRequest struct {
	BaseRequest
	Attributes LookupRequestAttributes `json:"attributes"`
}

type LookupRequestAttributes struct {
	Domain  string `json:"domain"`
	NoCache Bool   `json:"no_cache,omitempty"`
}

type LookupResponse struct {
	BaseResponse
	Attributes LookupResponseAttributes `json:"attributes"`
}

type LookupResponseAttributes struct {
	EmailAvailable Bool   `json:"email_available"`
	HasClaim       Bool   `json:"has_claim"`
	NoService      Bool   `json:"noservice"`
	PriceStatus    string `json:"price_status"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
}

type DomainsService struct {
	Client *Client
}

func (s *DomainsService) Lookup(attr LookupRequestAttributes) (*LookupResponse, error) {
	opsResponse := LookupResponse{}

	payload := LookupRequest{
		BaseRequest: BaseRequest{
			Action:   "LOOKUP",
			Object:   "DOMAIN",
			Protocol: "XCP",
		},
		Attributes: attr,
	}
	req, err := s.Client.NewRequest("POST", "", payload)
	if err != nil {
		return nil, err
	}
	err = s.Client.Do(req, &opsResponse)

	if err != nil {
		return nil, err
	}

	return &opsResponse, nil
}
