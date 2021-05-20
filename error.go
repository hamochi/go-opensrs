package opensrs

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Err             error
	HttpResponse    *http.Response
	OpenSRSResponse *BaseResponse
}

func (e ErrorResponse) Error() string {
	msg := ""
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if e.HttpResponse != nil {
		msg = msg + fmt.Sprintf("%s %s", e.HttpResponse.Request.Method, e.HttpResponse.Request.URL)
	}

	if e.OpenSRSResponse != nil {
		msg = msg + fmt.Sprintf("%s %s", e.OpenSRSResponse.ResponseCode, e.OpenSRSResponse.ResponseText)
	}
	return msg
}
