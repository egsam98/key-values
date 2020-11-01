package api

// swagger:response HTTPErrorResponse
type HTTPErrorResponse struct {
	// in: body
	Body struct {
		Message string `json:"message"`
	}
}

// swagger:response ValueResponse
type ValueResponse struct {
	// in: body
	Value Value
}

func NewValueResponse(value *string) *ValueResponse {
	return &ValueResponse{
		Value: Value{
			Value: value,
		},
	}
}
