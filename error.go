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
		msg = e.Err.Error() + "\n"
	}

	if e.HttpResponse != nil {
		msg = msg + "httpResponse: " + fmt.Sprintf("%s %s", e.HttpResponse.Request.Method, e.HttpResponse.Request.URL) + "\n"
	}

	if e.OpenSRSResponse != nil {
		msg = msg + "openSRSResponse " + fmt.Sprintf("%s %s", e.OpenSRSResponse.ResponseCode, e.OpenSRSResponse.ResponseText) + "\n"
	}
	return msg
}
