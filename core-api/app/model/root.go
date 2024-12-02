package model

type Response struct {
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	Success   bool        `json:"success"`
	ErrorCode *string     `json:"error_code"`
}
