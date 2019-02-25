package errs

// Failure struct implode runtime & api errors
type Failure struct {
	runtimeErr error
	apiErr     string
	apiErrs    map[string]string
}

// FailureResponse convert json error response into object
type FailureResponse struct {
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Errors   interface{} `json:"errors,omitempty"`
}
