package config

const FailedCode = 1
const SuccessCode = 0

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
