package dispatcher

type Request struct {
	Method    string        `json:"method"`
	Arguments []interface{} `json:"arguments"`

	// context key is ignored
}

type Response struct {
	Result interface{}    `json:"result"`
	Error  *ResponseError `json:"error"`

	Log string `json:"log"`
}

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`

	CanRetry bool `json:"ok_to_retry"`
}
