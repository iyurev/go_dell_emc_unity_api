package unity_api

import "fmt"

type RestErr struct {
	RespBody   []byte
	StatusCode int
}

func ok_status_code(status_code int) bool {
	success_status_codes := make(map[int]bool)
	success_status_codes[200] = true
	success_status_codes[202] = true
	success_status_codes[204] = true
	if _, ok := success_status_codes[status_code]; ok {
		return true
	}
	return false
}
func NewRestErr(b []byte, s int) *RestErr {
	return &RestErr{RespBody: b, StatusCode: s}
}

func (re *RestErr) Error() string {
	return fmt.Sprintf("Unity REST API ERROR!! Response code:  %s, Response body: %s\n", re.StatusCode, re.RespBody)
}
