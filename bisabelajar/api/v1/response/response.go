package response

type ResponseBody struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewSuccessResponse(data interface{}) ResponseBody {
	return ResponseBody{
		Status: "success",
		Data:   data,
	}
}

func NewErrorResponse(message string, err error) ResponseBody {
	return ResponseBody{
		Status:  "error",
		Message: message,
		Error:   err.Error(),
	}
}
