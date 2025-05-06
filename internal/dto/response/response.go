package response

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *Response) SetMessage(message string) {
	r.Message = message
}
