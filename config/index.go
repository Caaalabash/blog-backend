package config

const FailedCode = 0
const SuccessCode = 1

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
