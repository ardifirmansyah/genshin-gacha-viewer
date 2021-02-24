package api

import (
	"encoding/json"
	"time"

	"github.com/ardifirmansyah/genshin-gacha-viewer/src/common/constant"
)

type GeneralResponse struct {
	Success     bool           `json:"success"`
	ProcessTime string         `json:"process_time,omitempty"`
	Data        interface{}    `json:"data,omitempty"`
	Errors      []GeneralError `json:"errors,omitempty"`

	StartTime time.Time `json:"-"`
}

type GeneralError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewResponse() *GeneralResponse {
	return &GeneralResponse{
		StartTime: time.Now(),
	}
}

func (resp *GeneralResponse) AddError(err error) {
	code := constant.UnknownError

	resp.Errors = append(resp.Errors, GeneralError{
		Code:    code,
		Message: err.Error(),
	})
}

func (resp *GeneralResponse) AddErrorString(err string, code string) {
	resp.Errors = append(resp.Errors, GeneralError{
		Code:    code,
		Message: err,
	})
}

func (resp *GeneralResponse) JSON() []byte {
	if !resp.StartTime.IsZero() {
		resp.ProcessTime = time.Since(resp.StartTime).String()
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		return []byte(`
		{
			"success": false,
			"errors": [
				{
					"code": "500",
					"message": "Internal Server Error"
				}
			]
		}
	`)
	}

	return respJSON
}

