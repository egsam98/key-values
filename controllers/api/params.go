package api

// swagger:parameters getKeyValue
type GetKeyValueParams struct {
	// in: path
	// required: true
	Key string `json:"key"`
}

// swagger:parameters putKeyValue
type PutKeyValueParams struct {
	// in: path
	// required: true
	Key string `json:"key"`

	// in: body
	// required: true
	Kalue Value `json:"value"`
}
