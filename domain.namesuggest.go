package opensrs

// Response
type NameSuggestResponse struct {
	BaseResponse
	IsSearchComplete Bool                          `json:"is_search_complete"`
	SearchKey        string                        `json:"search_key"`
	Attributes       NameSuggestResponseAttributes `json:"attributes"`
}

type NameSuggestResponseAttributes struct {
	Lookup                  NameSuggestItems `json:"lookup"`
	PersonalNames           NameSuggestItems `json:"personal_names"`
	Premium                 NameSuggestItems `json:"premium"`
	PremiumBrokeredTransfer NameSuggestItems `json:"premium_brokered_transfer"`
	PremiumMakeOffer        NameSuggestItems `json:"premium_make_offer"`
	Suggestion              NameSuggestItems `json:"suggestion"`
}

type NameSuggestItems struct {
	Count        string            `json:"count"`
	ResponseText string            `json:"response_text"`
	ResponseCode string            `json:"response_code"`
	IsSuccess    Bool              `json:"is_success"`
	Items        []NameSuggestItem `json:"items"`
}

type NameSuggestItem struct {
	Domain             string `json:"domain"`
	PremiumPrice       string `json:"price"`
	Status             string `json:"status"`
	HasClaim           Bool   `json:"has_claim"`
	Reason             string `json:"reason"`
	ThirdPartyOfferUrl string `json:"third_party_offer_url"`
}

// Requests
type NameSuggestRequest struct {
	BaseRequest
	Attributes NameSuggestRequestAttributes `json:"attributes"`
}

type NameSuggestRequestAttributes struct {
	Languages          []string                   `json:"languages,omitempty"`
	MaxWaitTime        int                        `json:"max_wait_time,omitempty"`
	SearchKey          string                     `json:"search_key,omitempty"`
	SearchString       string                     `json:"searchstring,omitempty"`
	ServiceOverride    NameSuggestServiceOverride `json:"service_override,omitempty"`
	Services           []string                   `json:"services,omitempty"`
	SkipRegistryLookup bool                       `json:"skip_registry_lookup,omitempty"`
	TLDs               []string                   `json:"tlds,omitempty"`
}

type NameSuggestServiceOverride struct {
	Lookup        NameSuggestLookup     `json:"lookup,omitempty"`
	PersonalNames []string              `json:"personal_names,omitempty"`
	Premium       []string              `json:"premium,omitempty"`
	Suggestion    NameSuggestSuggestion `json:"suggestion,omitempty"`
}

type NameSuggestSuggestion struct {
	Maximum  string   `json:"maximum,omitempty"`
	PriceMax string   `json:"price_max,omitempty"`
	PriceMin string   `json:"price_min,omitempty"`
	TLDs     []string `json:"tlds,omitempty"`
}

type NameSuggestLookup struct {
	NameSuggestSuggestion
	NoCacheTld []string `json:"no_cache_tlds,omitempty"`
}

func (s *DomainsService) NameSuggest(attr NameSuggestRequestAttributes) (*NameSuggestResponse, error) {
	opsResponse := NameSuggestResponse{}

	payload := NameSuggestRequest{
		BaseRequest: BaseRequest{
			Action:   "NAME_SUGGEST",
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
