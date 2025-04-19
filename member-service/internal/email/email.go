package email

import (
	"encoding/json"
	"errors"
	"fmt"

	mlog "github.com/mike504110403/goutils/log"
	"github.com/valyala/fasthttp"
)

var SMTPURL string

func Init(url string) {
	SMTPURL = url
}

func SendEmail(bodyReq EmailRequest) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(SMTPURL)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	body, err := json.Marshal(bodyReq)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error marshalling request: %v", err))
		return err
	}
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	if err != nil {
		mlog.Error(fmt.Sprintf("Error forwarding request: %v", err))
		return err
	} else {
		if resp.StatusCode() != fasthttp.StatusOK {
			mlog.Error(fmt.Sprintf("Error forwarding request: %v", resp.StatusCode()))
			return errors.New("SMTP server error")
		}
	}

	return nil
}
