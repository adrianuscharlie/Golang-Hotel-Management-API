package response

type Response struct {
	ResponseCode    string `json:"responseCode" binding:"required"`
	ResponseMessage string `json:"responseMessage" binding:"required"`
	Data            any    `json:"data,omitempty"`
}
