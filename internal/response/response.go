package response

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
	}
}

// NewSuccessResponse creates a new success response with data
func NewSuccessResponse(code int, msg string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}
