package pbsdk

import (
	"fmt"

	prettygo "github.com/imsgao/pretty-go"
	"gopkg.in/resty.v1"
)

type Response struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Response) Error() error {
	if r.Code >= 200 && r.Code < 300 {
		return nil
	} else {
		return fmt.Errorf("code: %d, message: %s", r.Code, r.Message)
	}
}

func (r *Response) Status() string {
	return fmt.Sprintf("code: %d, message: %s", r.Code, r.Message)
}

func ResponseFromError(err error) *Response {
	return &Response{
		Code:    0,
		Message: err.Error(),
	}
}

func ResponseFromResty(resp *resty.Response) *Response {
	data := &Response{}
	if err := prettygo.DecodeJSON(resp.Body(), data); err != nil {
		return &Response{
			Code:    resp.StatusCode(),
			Message: err.Error(),
		}
	}
	return data
}

func ResponseWithStatus(code int, err error) *Response {
	return &Response{
		Code:    code,
		Message: err.Error(),
	}
}
