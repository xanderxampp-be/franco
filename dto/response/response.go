package response

type Response struct {
	Code   string      `json:"responseCode"`
	Refnum string      `json:"responseRefnum"`
	ID     string      `json:"responseId"`
	Desc   string      `json:"responseDesc"`
	Data   interface{} `json:"responseData"`
}

// Initialization of Response
func New(id string) *Response {
	return &Response{
		ID:   id,
		Code: "XX",
		Desc: "General Error",
		Data: new(struct{}),
	}
}
