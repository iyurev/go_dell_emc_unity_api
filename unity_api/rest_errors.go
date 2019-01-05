package unity_api

import "fmt"

type RestErr struct {
	RespBody   []byte
	StatusCode int
}

func OKStatusCode(status_code int) bool {
	success_codes := []int{200, 202, 204}
	for _, ok := range success_codes {
		if status_code == ok {
			return true
		}
	}
	return false
}

func NewRestErr(b []byte, s int) *RestErr {
	return &RestErr{RespBody: b, StatusCode: s}
}

func (re *RestErr) Error() string {
	return fmt.Sprintf("Unity REST API ERROR!! Response code:  %d, Response body: %s\n", re.StatusCode, re.RespBody)
}
