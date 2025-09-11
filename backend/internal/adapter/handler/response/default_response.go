package response

type DefaultResponse struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Message any  `json:"message"`
	Data    any  `json:"data"`
}

func ResponseError(code int, message any) DefaultResponse {
	return DefaultResponse{Success: false, Code: code, Message: message, Data: nil}
}

func ResponseSuccess(code int, message, data any) DefaultResponse {
	return DefaultResponse{Success: true, Code: code, Message: message, Data: data}
}
